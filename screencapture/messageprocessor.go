package screencapture

import (
	"encoding/binary"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	log "github.com/sirupsen/logrus"
)

//MessageProcessor is used to implement the control flow of USB messages answers and replies.
//It receives readily split byte frames, parses them, responds to them and passes on
//extracted CMSampleBuffers to a consumer
type MessageProcessor struct {
	usbWriter                                UsbWriter
	stopSignal                               chan interface{}
	clock                                    coremedia.CMClock
	localAudioClock                          coremedia.CMClock
	needClockRef                             packet.CFTypeID
	needMessage                              []byte
	audioSamplesReceived                     int
	videoSamplesReceived                     int
	cmSampleBufConsumer                      CmSampleBufConsumer
	clockBuilder                             func(uint64) coremedia.CMClock
	deviceAudioClockRef                      packet.CFTypeID
	releaseWaiter                            chan interface{}
	firstAudioTimeTaken                      bool
	startTimeDeviceAudioClock                coremedia.CMTime
	startTimeLocalAudioClock                 coremedia.CMTime
	lastEatFrameReceivedDeviceAudioClockTime coremedia.CMTime
	lastEatFrameReceivedLocalAudioClockTime  coremedia.CMTime
}

//NewMessageProcessor creates a new MessageProcessor that will write answers to the given UsbWriter,
// forward extracted CMSampleBuffers to the CMSampleBufConsumer and wait for the stopSignal.
func NewMessageProcessor(usbWriter UsbWriter, stopSignal chan interface{}, consumer CmSampleBufConsumer) MessageProcessor {
	clockBuilder := func(ID uint64) coremedia.CMClock { return coremedia.NewCMClockWithHostTime(ID) }
	return NewMessageProcessorWithClockBuilder(usbWriter, stopSignal, consumer, clockBuilder)
}

//NewMessageProcessorWithClockBuilder lets you inject a clockBuilder for the sake of testability.
func NewMessageProcessorWithClockBuilder(usbWriter UsbWriter, stopSignal chan interface{}, consumer CmSampleBufConsumer, clockBuilder func(uint64) coremedia.CMClock) MessageProcessor {
	var mp = MessageProcessor{usbWriter: usbWriter, stopSignal: stopSignal, cmSampleBufConsumer: consumer, clockBuilder: clockBuilder, releaseWaiter: make(chan interface{}), firstAudioTimeTaken: false}
	return mp
}

//ReceiveData waits for byte frames of the correct length without the length field.
//This function will only accept byte frames starting with the ASYN, SYNC or PING uint32 magic.
func (mp *MessageProcessor) ReceiveData(data []byte) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.PingPacketMagic:
		log.Info("AudioVideo-Stream has started")
		mp.usbWriter.WriteDataToUsb(packet.NewPingPacketAsBytes())
		return
	case packet.SyncPacketMagic:
		mp.handleSyncPacket(data)
		return
	case packet.AsynPacketMagic:
		mp.handleAsyncPacket(data)
		return
	default:
		log.Warnf("received unknown packet type: %x", data[:4])
	}
	mp.stop()
}

func (mp *MessageProcessor) handleSyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.OG:
		ogPacket, err := packet.NewSyncOgPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing OG packet", err)
		}
		log.Debugf("Rcv:%s", ogPacket.String())

		replyBytes := ogPacket.NewReply()
		log.Debugf("Send OG-REPLY {correlation:%x}", ogPacket.CorrelationID)
		mp.usbWriter.WriteDataToUsb(replyBytes)
	case packet.CWPA:
		cwpaPacket, err := packet.NewSyncCwpaPacketFromBytes(data)
		if err != nil {
			log.Error("failed parsing cwpa packet", err)
			return
		}
		log.Debugf("Rcv:%s", cwpaPacket.String())
		clockRef := cwpaPacket.DeviceClockRef + 1000

		mp.localAudioClock = coremedia.NewCMClockWithHostTime(clockRef)
		mp.deviceAudioClockRef = cwpaPacket.DeviceClockRef
		deviceInfo := packet.NewAsynHpd1Packet(packet.CreateHpd1DeviceInfoDict())
		log.Debug("Sending ASYN HPD1")
		mp.usbWriter.WriteDataToUsb(deviceInfo)
		log.Debugf("Send CWPA-RPLY {correlation:%x, clockRef:%x}", cwpaPacket.CorrelationID, clockRef)
		mp.usbWriter.WriteDataToUsb(cwpaPacket.NewReply(clockRef))
		log.Debug("Sending ASYN HPD1")
		mp.usbWriter.WriteDataToUsb(deviceInfo)
		deviceInfo1 := packet.NewAsynHpa1Packet(packet.CreateHpa1DeviceInfoDict(), cwpaPacket.DeviceClockRef)
		log.Debug("Sending ASYN HPA1")
		mp.usbWriter.WriteDataToUsb(deviceInfo1)
	case packet.CVRP:
		cvrpPacket, err := packet.NewSyncCvrpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing CVRP packet", err)
			return
		}

		log.Debugf("Rcv:%s", cvrpPacket.String())

		mp.needClockRef = cvrpPacket.DeviceClockRef
		mp.needMessage = packet.AsynNeedPacketBytes(mp.needClockRef)
		log.Debugf("Send NEED %x", mp.needClockRef)
		mp.usbWriter.WriteDataToUsb(mp.needMessage)

		clockRef2 := cvrpPacket.DeviceClockRef + 0x1000AF
		log.Debugf("Send CVRP-RPLY {correlation:%x, clockRef:%x}", cvrpPacket.CorrelationID, clockRef2)
		mp.usbWriter.WriteDataToUsb(cvrpPacket.NewReply(clockRef2))
	case packet.CLOK:
		clokPacket, err := packet.NewSyncClokPacketFromBytes(data)
		if err != nil {
			log.Error("Failed parsing Clok Packet", err)
		}
		log.Debugf("Rcv:%s", clokPacket.String())
		clockRef := clokPacket.ClockRef + 0x10000
		mp.clock = coremedia.NewCMClockWithHostTime(clockRef)
		log.Debugf("Send CLOK-RPLY {correlation:%x, clockRef:%x}", clokPacket.CorrelationID, clockRef)
		mp.usbWriter.WriteDataToUsb(clokPacket.NewReply(clockRef))
	case packet.TIME:
		timePacket, err := packet.NewSyncTimePacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing TIME SYNC packet", err)
		}
		log.Debugf("Rcv:%s", timePacket.String())
		timeToSend := mp.clock.GetTime()
		replyBytes, err := timePacket.NewReply(timeToSend)
		if err != nil {
			log.Error("Could not create SYNC TIME REPLY")
		}
		log.Debugf("Send TIME-REPLY {correlation:%x, time:%s}", timePacket.CorrelationID, timeToSend)
		mp.usbWriter.WriteDataToUsb(replyBytes)
	case packet.AFMT:
		afmtPacket, err := packet.NewSyncAfmtPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC AFMT packet", err)
		}
		log.Debugf("Rcv:%s", afmtPacket.String())

		replyBytes := afmtPacket.NewReply()
		log.Debugf("Send AFMT-REPLY {correlation:%x}", afmtPacket.CorrelationID)
		mp.usbWriter.WriteDataToUsb(replyBytes)
	case packet.SKEW:
		skewPacket, err := packet.NewSyncSkewPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC SKEW packet", err)
		}
		skewValue := coremedia.CalculateSkew(mp.startTimeLocalAudioClock, mp.lastEatFrameReceivedLocalAudioClockTime, mp.startTimeDeviceAudioClock, mp.lastEatFrameReceivedDeviceAudioClockTime)
		log.Debugf("Rcv:%s Reply:%f", skewPacket.String(), skewValue)
		mp.usbWriter.WriteDataToUsb(skewPacket.NewReply(skewValue))
	case packet.STOP:
		stopPacket, err := packet.NewSyncStopPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC STOP packet", err)
		}
		log.Debugf("Rcv:%s", stopPacket.String())
		mp.usbWriter.WriteDataToUsb(stopPacket.NewReply())
	default:
		log.Warnf("received unknown sync packet type: %x", data)
		mp.stop()
	}
}

func (mp *MessageProcessor) handleAsyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.EAT:
		mp.audioSamplesReceived++
		eatPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(data)
		if err != nil {
			log.Warn("unknown eat")
			return
		}
		if !mp.firstAudioTimeTaken {
			mp.startTimeDeviceAudioClock = eatPacket.CMSampleBuf.OutputPresentationTimestamp
			mp.startTimeLocalAudioClock = mp.localAudioClock.GetTime()
			mp.lastEatFrameReceivedDeviceAudioClockTime = eatPacket.CMSampleBuf.OutputPresentationTimestamp
			mp.lastEatFrameReceivedLocalAudioClockTime = mp.startTimeLocalAudioClock
			mp.firstAudioTimeTaken = true
		} else {
			mp.lastEatFrameReceivedDeviceAudioClockTime = eatPacket.CMSampleBuf.OutputPresentationTimestamp
			mp.lastEatFrameReceivedLocalAudioClockTime = mp.localAudioClock.GetTime()
		}

		err = mp.cmSampleBufConsumer.Consume(eatPacket.CMSampleBuf)
		if err != nil {
			log.Warn("failed consuming audio buf", err)
			return
		}
		if mp.audioSamplesReceived%100 == 0 {
			log.Debugf("RCV Audio Samples:%d", mp.audioSamplesReceived)
		}
	case packet.FEED:
		feedPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(data)
		if err != nil {
			log.Errorf("Error parsing FEED packet: %x %s", data, err)
			mp.usbWriter.WriteDataToUsb(mp.needMessage)
			return
		}
		mp.videoSamplesReceived++
		err = mp.cmSampleBufConsumer.Consume(feedPacket.CMSampleBuf)
		if err != nil {
			log.Fatal("Failed writing sample data to Consumer", err)
		}
		if mp.videoSamplesReceived%500 == 0 {
			log.Debugf("Rcv'd(%d) last:%s", mp.videoSamplesReceived, feedPacket.String())
			mp.videoSamplesReceived = 0
		}

		mp.usbWriter.WriteDataToUsb(mp.needMessage)
	case packet.SPRP:
		sprpPacket, err := packet.NewAsynSprpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SPRP packet", err)
			return
		}
		log.Debugf("Rcv:%s", sprpPacket.String())
	case packet.TJMP:
		tjmpPacket, err := packet.NewAsynTjmpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing tjmp packet", err)
			return
		}
		log.Debugf("Rcv:%s", tjmpPacket.String())
	case packet.SRAT:
		sratPacket, err := packet.NewAsynSratPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing srat packet", err)
			return
		}
		log.Debugf("Rcv:%s", sratPacket.String())
	case packet.TBAS:
		tbasPacket, err := packet.NewAsynTbasPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing tbas packet", err)
			return
		}
		log.Debugf("Rcv:%s", tbasPacket.String())
	case packet.RELS:
		relsPacket, err := packet.NewAsynRelsPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing RELS packet", err)
			return
		}
		log.Debugf("Rcv:%s", relsPacket.String())
		var signal interface{}
		mp.releaseWaiter <- signal
	default:
		log.Warnf("received unknown async packet type: %x", data)
		mp.stop()
	}
}

//CloseSession shuts down the streams on the device by sending HPA0 and HPD0
//messages and waiting for RELS messages.
func (mp *MessageProcessor) CloseSession() {
	log.Info("Telling device to stop streaming..")
	mp.usbWriter.WriteDataToUsb(packet.NewAsynHPA0(mp.deviceAudioClockRef))
	mp.usbWriter.WriteDataToUsb(packet.NewAsynHPD0())
	log.Info("Waiting for device to tell us to stop..")
	for i := 0; i < 2; i++ {
		select {
		case <-mp.releaseWaiter:
		case <-time.After(3 * time.Second):
			log.Warn("Timed out waiting for device closing")
			return
		}
	}
	mp.usbWriter.WriteDataToUsb(packet.NewAsynHPD0())
	log.Info("OK. Ready to release USB Device.")
}

func (mp MessageProcessor) stop() {
	var stop interface{}
	mp.stopSignal <- stop
}
