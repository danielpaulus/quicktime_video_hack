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

func NewMessageProcessor(writeToUsb func([]byte), stopSignal chan interface{}) messageProcessor {
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
		mp.clock = coremedia.CMClock{TimeScale: 1000000000, ID: clockRef}
		log.Debug("Sending CLOK reply")
		mp.writeToUsb(clok.NewReply(clockRef))
	case packet.TIME:
		log.Debug("Received Sync Time")
		timePacket, err := packet.NewSyncTimePacketFromBytes(data)
		if err != nil {
			log.Error("Error parsing TIME SYNC packet", err)
		}
		replyBytes, err := timePacket.NewReply(mp.clock.GetTime())
		if err != nil {
			log.Error("Could not create SYNC TIME REPLY")
		}
		log.Debug("Sending TIME REPLY")
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
	default:
		log.Warnf("received unknown async packet type: %x", data)
	}
}
