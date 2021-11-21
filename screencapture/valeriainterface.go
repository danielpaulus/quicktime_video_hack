package screencapture

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	log "github.com/sirupsen/logrus"
	"time"
)

type ValeriaInterface struct {
	Local        LocalValeriaApi
	Remote       DeviceValeriaAPI
	errorChannel chan error
	closeChannel chan interface{}
}

type DataHolder struct {
	localAudioClock                          coremedia.CMClock
	deviceAudioClockRef                      packet.CFTypeID
	needClockRef                             packet.CFTypeID
	clock                                    coremedia.CMClock
	startTimeLocalAudioClock                 coremedia.CMTime
	lastEatFrameReceivedLocalAudioClockTime  coremedia.CMTime
	startTimeDeviceAudioClock                coremedia.CMTime
	lastEatFrameReceivedDeviceAudioClockTime coremedia.CMTime
	audioSamplesReceived                     uint64
	firstAudioTimeTaken                      bool
	videoSamplesReceived                     uint64
}

type LocalValeriaApi struct {
	remote              DeviceValeriaAPI
	pingChannel         chan error
	audioClockChannel   chan error
	dataHolder          DataHolder
	videoClockChannel   chan error
	sampleDataChannel   chan coremedia.CMSampleBuffer
	audioReleaseChannel chan error
	videoReleaseChannel chan error
}

type DeviceValeriaAPI struct {
	usbAdapter *UsbAdapterNew
	dataHolder DataHolder
}

func NewValeriaInterface() ValeriaInterface {
	dataHolder := DataHolder{}
	local := LocalValeriaApi{
		pingChannel:         make(chan error, 1),
		audioClockChannel:   make(chan error, 1),
		videoClockChannel:   make(chan error, 1),
		audioReleaseChannel: make(chan error, 1),
		videoReleaseChannel: make(chan error, 1),
		dataHolder:          dataHolder,
		sampleDataChannel:   make(chan coremedia.CMSampleBuffer, 50),
	}
	remote := DeviceValeriaAPI{dataHolder: dataHolder}
	valeriaIface := ValeriaInterface{Local: local,
		errorChannel: make(chan error, 1),
		closeChannel: make(chan interface{}),
		Remote:       remote,
	}
	return valeriaIface
}

func (l LocalValeriaApi) AwaitAudioClockRelease() error {
	return awaitOrTimeout(l.audioReleaseChannel, "audio clock release")
}

func (l LocalValeriaApi) AwaitVideoClockRelease() error {
	return awaitOrTimeout(l.videoReleaseChannel, "video clock release")
}

func (l LocalValeriaApi) AwaitVideoClockSync() error {
	return awaitOrTimeout(l.videoClockChannel, "video clock")
}

func (l LocalValeriaApi) AwaitAudioClockSync() error {
	return awaitOrTimeout(l.audioClockChannel, "audio clock")
}

func (l LocalValeriaApi) AwaitPing() error {
	return awaitOrTimeout(l.pingChannel, "waiting for ping")
}

func awaitOrTimeout(channel chan error, loggerTag string) error {
	select {
	case <-time.After(5 * time.Second):
		return fmt.Errorf("%s timed out. restart the device please it might be buggy", loggerTag)
	case err := <-channel:
		if err != nil {
			return fmt.Errorf("failed '%s' with device %v", loggerTag, err)
		}
	}
	log.Infof("%s succeeded", loggerTag)
	return nil
}

func (l LocalValeriaApi) ping() {
	l.pingChannel <- nil
}

//I don't know what the !go command is for. It seems we always just return 0 and it works.
func (l LocalValeriaApi) gocmd(unknown uint32) uint64 {
	log.Debugf("go! %d", unknown)
	return 0
}

func (l *LocalValeriaApi) setupAudioClock(deviceClockRef packet.CFTypeID) packet.CFTypeID {
	clockRef := deviceClockRef + 1000

	l.dataHolder.localAudioClock = coremedia.NewCMClockWithHostTime(clockRef)
	l.dataHolder.deviceAudioClockRef = deviceClockRef
	l.audioClockChannel <- nil
	return clockRef
}

//this message is always the same, so we just prepare it once and send the same bytes all the time
var needMessage []byte

func (l *LocalValeriaApi) setupVideoClock(deviceClockRef packet.CFTypeID) packet.CFTypeID {
	l.dataHolder.needClockRef = deviceClockRef
	needMessage = packet.AsynNeedPacketBytes(deviceClockRef)
	l.videoClockChannel <- nil
	return deviceClockRef + 0x1000AF
}

func (l LocalValeriaApi) setupMainClock(ref packet.CFTypeID) packet.CFTypeID {
	clockRef := ref + 0x10000
	l.dataHolder.clock = coremedia.NewCMClockWithHostTime(clockRef)
	return clockRef
}

func (l LocalValeriaApi) time() coremedia.CMTime {
	return l.dataHolder.clock.GetTime()
}

func (l LocalValeriaApi) stop() {
	log.Info("device sent STOP command")
}

func (l LocalValeriaApi) skew() float64 {
	return coremedia.CalculateSkew(
		l.dataHolder.startTimeLocalAudioClock,
		l.dataHolder.lastEatFrameReceivedLocalAudioClockTime,
		l.dataHolder.startTimeDeviceAudioClock,
		l.dataHolder.lastEatFrameReceivedDeviceAudioClockTime)
}

//ReadSampleBuffer blocks until a buffer is received or the interface is closed
func (l LocalValeriaApi) ReadSampleBuffer() coremedia.CMSampleBuffer {
	return <-l.sampleDataChannel
}

func (l LocalValeriaApi) receiveAudioSample(buf coremedia.CMSampleBuffer) {
	if !l.dataHolder.firstAudioTimeTaken {
		l.dataHolder.startTimeDeviceAudioClock = buf.OutputPresentationTimestamp
		l.dataHolder.startTimeLocalAudioClock = l.dataHolder.localAudioClock.GetTime()
		l.dataHolder.lastEatFrameReceivedDeviceAudioClockTime = buf.OutputPresentationTimestamp
		l.dataHolder.lastEatFrameReceivedLocalAudioClockTime = l.dataHolder.startTimeLocalAudioClock
		l.dataHolder.firstAudioTimeTaken = true
	} else {
		l.dataHolder.lastEatFrameReceivedDeviceAudioClockTime = buf.OutputPresentationTimestamp
		l.dataHolder.lastEatFrameReceivedLocalAudioClockTime = l.dataHolder.localAudioClock.GetTime()
	}

	l.sampleDataChannel <- buf
	if log.IsLevelEnabled(log.DebugLevel) {
		l.dataHolder.audioSamplesReceived++
		if l.dataHolder.audioSamplesReceived%100 == 0 {
			log.Debugf("RCV Audio Samples:%d", l.dataHolder.audioSamplesReceived)
		}
	}
}

func (l LocalValeriaApi) feed(buf coremedia.CMSampleBuffer) {
	l.sampleDataChannel <- buf
	if log.IsLevelEnabled(log.DebugLevel) {
		l.dataHolder.videoSamplesReceived++
		if l.dataHolder.videoSamplesReceived%500 == 0 {
			log.Debugf("Rcv'd(%d) last:%s", l.dataHolder.videoSamplesReceived, buf.String())
			l.dataHolder.videoSamplesReceived = 0
		}
	}
}

func (l LocalValeriaApi) timeJump(unknown []byte) {

}

func (l LocalValeriaApi) setProperties(property coremedia.StringKeyEntry, clockRef packet.CFTypeID) {

}

func (l LocalValeriaApi) setClockRate(clockRef packet.CFTypeID, cmTime coremedia.CMTime, rate1 float32, rate2 float32) {

}

func (l LocalValeriaApi) setTimeBase(ref packet.CFTypeID, ref2 packet.CFTypeID) {

}

func (l LocalValeriaApi) release(clockRef packet.CFTypeID) {
	if clockRef == l.dataHolder.needClockRef {
		l.videoReleaseChannel <- nil
		return
	}
	if clockRef == l.dataHolder.localAudioClock.ID {
		l.audioReleaseChannel <- nil
		return
	}
	log.Warnf("release for unknown clock received %d", clockRef)
}

func (d DeviceValeriaAPI) RequestSampleData() error {
	log.Debugf("Send NEED %x", d.dataHolder.needClockRef)
	return d.usbAdapter.WriteDataToUsb(needMessage)
}

func (d DeviceValeriaAPI) EnableVideo() error {
	deviceInfo := packet.NewAsynHpd1Packet(packet.CreateHpd1DeviceInfoDict())
	log.Debug("Sending ASYN HPD1")
	err := d.usbAdapter.WriteDataToUsb(deviceInfo)
	if err != nil {
		return err
	}
	log.Debug("Sending ASYN HPD1")
	return d.usbAdapter.WriteDataToUsb(deviceInfo)
}

func (d DeviceValeriaAPI) EnableAudio() error {

	deviceInfo1 := packet.NewAsynHpa1Packet(packet.CreateHpa1DeviceInfoDict(), d.dataHolder.deviceAudioClockRef)
	log.Debug("Sending ASYN HPA1")
	return d.usbAdapter.WriteDataToUsb(deviceInfo1)
}

func (d DeviceValeriaAPI) Ping() error {
	return d.usbAdapter.WriteDataToUsb(packet.NewPingPacketAsBytes())
}

// StartReadLoop claims&opens the USB Device and starts listening to RPC calls
// and blocks until ValeriaInterface is closed or an error occurs.
func (v *ValeriaInterface) StartReadLoop(device IosDevice) error {
	usbAdapter := &UsbAdapterNew{}
	err := usbAdapter.InitializeUSB(device)
	if err != nil {
		return fmt.Errorf("failed initializing usb with error %v", err)
	}
	v.Remote.usbAdapter = usbAdapter
	return readLoop(v, usbAdapter)
}

//readLoop reads messages sent by the device and dispatches them to the local api
func readLoop(v *ValeriaInterface, usbAdapter *UsbAdapterNew) error {
	for {
		select {
		case err := <-v.errorChannel:
			return err
		case <-v.closeChannel:
			return nil
		default:
			frame, err := usbAdapter.ReadFrame()
			if err != nil {
				return err
			}
			handleFrame(frame, v, usbAdapter)
		}

	}
}

// Decode Remote rpc calls and forward them to the local Valeria API
func handleFrame(data []byte, valeria *ValeriaInterface, usbAdapter *UsbAdapterNew) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.PingPacketMagic:
		valeria.Local.ping()
	case packet.SyncPacketMagic:
		err := handleSyncPacket(data, valeria, usbAdapter)
		if err != nil {
			valeria.errorChannel <- err
			return
		}
		return
	case packet.AsynPacketMagic:
		err := handleAsyncPacket(data, valeria)
		if err != nil {
			valeria.errorChannel <- err
			return
		}
		return
	default:
		valeria.errorChannel <- fmt.Errorf("received unknown packet type: %x", data[:4])
	}
}

func handleSyncPacket(data []byte, valeria *ValeriaInterface, adapter *UsbAdapterNew) error {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.OG:
		ogPacket, err := packet.NewSyncOgPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing OG packet", err)
		}
		log.Debugf("Rcv:%s", ogPacket.String())
		response := valeria.Local.gocmd(ogPacket.Unknown)
		replyBytes := ogPacket.NewReply(response)
		return adapter.WriteDataToUsb(replyBytes)
	case packet.CWPA:
		cwpaPacket, err := packet.NewSyncCwpaPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("failed parsing cwpa packet %v", err)
		}
		log.Debugf("Rcv:%s", cwpaPacket.String())
		clockRef := valeria.Local.setupAudioClock(cwpaPacket.DeviceClockRef)

		log.Debugf("Send CWPA-RPLY {correlation:%x, clockRef:%x}", cwpaPacket.CorrelationID, clockRef)
		return adapter.WriteDataToUsb(cwpaPacket.NewReply(clockRef))

	case packet.CVRP:
		cvrpPacket, err := packet.NewSyncCvrpPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing CVRP packet %v", err)
		}
		log.Debugf("Rcv:%s", cvrpPacket.String())
		videoClockRef := valeria.Local.setupVideoClock(cvrpPacket.DeviceClockRef)

		log.Debugf("Send CVRP-RPLY {correlation:%x, clockRef:%x}", cvrpPacket.CorrelationID, videoClockRef)
		return adapter.WriteDataToUsb(cvrpPacket.NewReply(videoClockRef))
	case packet.CLOK:
		clokPacket, err := packet.NewSyncClokPacketFromBytes(data)
		if err != nil {
			log.Error("Failed parsing Clok Packet", err)
		}
		log.Debugf("Rcv:%s", clokPacket.String())
		clockRef := valeria.Local.setupMainClock(clokPacket.ClockRef)

		log.Debugf("Send CLOK-RPLY {correlation:%x, clockRef:%x}", clokPacket.CorrelationID, clockRef)
		return adapter.WriteDataToUsb(clokPacket.NewReply(clockRef))
	case packet.TIME:
		timePacket, err := packet.NewSyncTimePacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing TIME SYNC packet", err)
		}
		log.Debugf("Rcv:%s", timePacket.String())
		timeToSend := valeria.Local.time()
		replyBytes, err := timePacket.NewReply(timeToSend)
		if err != nil {
			return fmt.Errorf("could not create SYNC TIME REPLY")
		}
		log.Debugf("Send TIME-REPLY {correlation:%x, time:%s}", timePacket.CorrelationID, timeToSend)
		return adapter.WriteDataToUsb(replyBytes)
		//TODO: turn into nice API function
	case packet.AFMT:
		afmtPacket, err := packet.NewSyncAfmtPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC AFMT packet", err)
		}
		log.Debugf("Rcv:%s", afmtPacket.String())

		replyBytes := afmtPacket.NewReply()
		log.Debugf("Send AFMT-REPLY {correlation:%x}", afmtPacket.CorrelationID)
		return adapter.WriteDataToUsb(replyBytes)
	case packet.SKEW:
		skewPacket, err := packet.NewSyncSkewPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC SKEW packet", err)
		}
		skewValue := valeria.Local.skew()
		log.Debugf("Rcv:%s Reply:%f", skewPacket.String(), skewValue)
		return adapter.WriteDataToUsb(skewPacket.NewReply(skewValue))
	case packet.STOP:
		stopPacket, err := packet.NewSyncStopPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SYNC STOP packet", err)
		}
		valeria.Local.stop()
		log.Debugf("Rcv:%s", stopPacket.String())
		return adapter.WriteDataToUsb(stopPacket.NewReply())
	default:
		return fmt.Errorf("received unknown sync packet type: %x", data)
	}
}

func handleAsyncPacket(data []byte, valeria *ValeriaInterface) error {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.EAT:
		eatPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("eat packet could not be unmarshalled %v", err)
		}
		valeria.Local.receiveAudioSample(eatPacket.CMSampleBuf)
		return nil
	case packet.FEED:
		feedPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing FEED packet: %x %s", data, err)
		}
		valeria.Local.feed(feedPacket.CMSampleBuf)
		return nil
	case packet.SPRP:
		sprpPacket, err := packet.NewAsynSprpPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing SPRP packet %v", err)
		}
		valeria.Local.setProperties(sprpPacket.Property, sprpPacket.ClockRef)
		log.Debugf("Rcv:%s", sprpPacket.String())
		return nil
	case packet.TJMP:
		tjmpPacket, err := packet.NewAsynTjmpPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing tjmp packet %v", err)
		}
		valeria.Local.timeJump(tjmpPacket.Unknown)
		log.Debugf("Rcv:%s", tjmpPacket.String())
		return nil
	case packet.SRAT:
		sratPacket, err := packet.NewAsynSratPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing srat packet %v", err)
		}
		valeria.Local.setClockRate(sratPacket.ClockRef, sratPacket.Time, sratPacket.Rate1, sratPacket.Rate2)
		log.Debugf("Rcv:%s", sratPacket.String())
		return nil
	case packet.TBAS:
		tbasPacket, err := packet.NewAsynTbasPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing tbas packet %v", err)
		}
		valeria.Local.setTimeBase(tbasPacket.ClockRef, tbasPacket.SomeOtherRef)
		log.Debugf("Rcv:%s", tbasPacket.String())
		return nil
	case packet.RELS:
		relsPacket, err := packet.NewAsynRelsPacketFromBytes(data)
		if err != nil {
			return fmt.Errorf("error parsing RELS packet %v", err)
		}
		valeria.Local.release(relsPacket.ClockRef)
		log.Debugf("Rcv:%s", relsPacket.String())
		return nil
	default:
		return fmt.Errorf("received unknown async packet type: %x", data)
	}
}
