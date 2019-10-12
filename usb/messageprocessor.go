package usb

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/usb/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	log "github.com/sirupsen/logrus"
)

type messageProcessor struct {
	connectionState    int
	writeToUsb         func([]byte)
	stopSignal         chan interface{}
	clock              coremedia.CMClock
	totalBytesReceived int
}

func newMessageProcessor(writeToUsb func([]byte), stopSignal chan interface{}) messageProcessor {
	var mp = messageProcessor{writeToUsb: writeToUsb, stopSignal: stopSignal}
	return mp
}

func (mp *messageProcessor) receiveData(data []byte) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.PingPacketMagic:
		log.Debug("initial ping received, sending ping back")
		mp.writeToUsb(packet.NewPingPacketAsBytes())
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

	var stop interface{}
	mp.stopSignal <- stop
}

func (mp *messageProcessor) handleSyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.CWPA:
		cwpaPacket, err := packet.NewSyncCwpaPacketFromBytes(data)
		if err != nil {
			log.Error("failed parsing cwpa packet", err)
			return
		}
		log.Debugf("Rcv:%s", cwpaPacket.String())
		clockRef := cwpaPacket.DeviceClockRef + 1000

		deviceInfo := packet.NewAsynHpd1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debug("Sending ASYN HPD1")
		mp.writeToUsb(deviceInfo)
		log.Debugf("Sending CWPA Reply:%x", clockRef)
		mp.writeToUsb(cwpaPacket.NewReply(clockRef))
		log.Debug("Sending ASYN HPD1")
		mp.writeToUsb(deviceInfo)
		deviceInfo1 := packet.NewAsynHpa1Packet(messages.CreateHpa1DeviceInfoDict(), clockRef)
		log.Debug("Sending ASYN HPA1")
		mp.writeToUsb(deviceInfo1)
	case packet.CVRP:
		cvrpPacket, err := packet.NewSyncCvrpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing CVRP packet", err)
			return
		}
		log.Debugf("Rcv:%s", cvrpPacket.String())
		clockRef2 := cvrpPacket.DeviceClockRef + 1000
		log.Debugf("Sending CVRP Reply:%x", clockRef2)
		mp.writeToUsb(cvrpPacket.NewReply(clockRef2))
	case packet.CLOK:
		clokPacket, err := packet.NewSyncClokPacketFromBytes(data)
		if err != nil {
			log.Error("Failed parsing Clok Packet", err)
		}
		log.Debugf("Rcv:%s", clokPacket.String())
		clockRef := clokPacket.ClockRef + 0x10000
		mp.clock = coremedia.NewCMClockWithHostTime(clockRef)
		log.Debugf("Sending CLOK reply:%x", clockRef)
		mp.writeToUsb(clokPacket.NewReply(clockRef))
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
		log.Debugf("Sending TIME REPLY:%s", timeToSend.String())
		mp.writeToUsb(replyBytes)
	default:
		log.Warnf("received unknown sync packet type: %x", data)
	}
}

func (mp *messageProcessor) handleAsyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data[12:]) {
	case packet.FEED:
		feedPacket, err := packet.NewAsynFeedPacketFromBytes(data)
		if err != nil {
			log.Errorf("Error parsing FEED packet: %x", data, err)
			return
		}
		log.Debugf("Rcv:%s", feedPacket.String())
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
	default:
		log.Warnf("received unknown async packet type: %x", data)
	}
}
