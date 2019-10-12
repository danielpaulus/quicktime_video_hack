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
		log.Debug("Received Sync CWPA")
		cwpaPacket, err := packet.NewSyncCwpaPacketFromBytes(data)
		if err != nil {
			log.Error("failed parsing cwpa packet", err)
			return
		}
		clockRef := cwpaPacket.DeviceClockRef + 1000

		deviceInfo := packet.NewAsynHpd1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debug("Sending ASYN HPD1")
		mp.writeToUsb(deviceInfo)
		log.Debug("Sending CWPA Reply")
		mp.writeToUsb(cwpaPacket.NewReply(clockRef))
		log.Debug("Sending ASYN HPD1")
		mp.writeToUsb(deviceInfo)
		deviceInfo1 := packet.NewAsynHpa1Packet(messages.CreateHpa1DeviceInfoDict(), clockRef)
		log.Debug("Sending ASYN HPA1")
		mp.writeToUsb(deviceInfo1)
	case packet.CVRP:
		log.Debug("Received Sync CVRP")
		cvrpPacket, err := packet.NewSyncCvrpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing CVRP packet", err)
			return
		}
		clockRef2 := cvrpPacket.DeviceClockRef + 1000
		log.Debug("Sending CVRP Reply")
		mp.writeToUsb(cvrpPacket.NewReply(clockRef2))
		log.Debugf("CVRP:%s", cvrpPacket.Payload.String())
	case packet.CLOK:
		log.Debug("Received Sync Clock")
		clok, err := packet.NewSyncClokPacketFromBytes(data)
		if err != nil {
			log.Error("Failed parsing Clok Packet", err)
		}
		clockRef := clok.ClockRef + 0x10000
		mp.clock = coremedia.NewCMClockWithHostTime(clockRef)
		log.Debug("Sending CLOK reply")
		mp.writeToUsb(clok.NewReply(clockRef))
	case packet.TIME:
		log.Debug("Received Sync Time")
		timePacket, err := packet.NewSyncTimePacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing TIME SYNC packet", err)
		}
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
		mp.totalBytesReceived += len(data)
		log.Debugf("rcv feed: %d bytes - %d total", len(data), mp.totalBytesReceived)
	//mp.writeToUsb(packet.AsynNeedPacketBytes)
	case packet.SPRP:
		packet, err := packet.NewAsynSprpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing SPRP packet", err)
			return
		}
		log.Debugf("rcv set property (sprp):%s", packet.Property.Key)
	case packet.TJMP:
		packet, err := packet.NewAsynTjmpPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing tjmp packet", err)
			return
		}
		log.Debugf("rcv tjmp: 0x%x", packet.Unknown)
	case packet.SRAT:
		packet, err := packet.NewAsynSratPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing srat packet", err)
			return
		}
		log.Debugf("rcv srat: rate1:%f rate2:%f time:%s", packet.Rate1, packet.Rate2, packet.Time.String())
	case packet.TBAS:
		packet, err := packet.NewAsynTbasPacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing tbas packet", err)
			return
		}
		log.Debugf("rcv tbas: 0x%x", packet.SomeOtherRef)
	default:
		log.Warnf("received unknown async packet type: %x", data)
	}
}
