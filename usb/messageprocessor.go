package usb

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	log "github.com/sirupsen/logrus"
)

const (
	initialState          = iota
	pingSent              = iota
	pingExchangeCompleted = iota
)

type messageProcessor struct {
	connectionState int
	writeToUsb      func([]byte)
	stopSignal      chan interface{}
}

func NewMessageProcessor(writeToUsb func([]byte), stopSignal chan interface{}) messageProcessor {
	var mp = messageProcessor{connectionState: initialState, writeToUsb: writeToUsb, stopSignal: stopSignal}
	return mp
}

func (mp *messageProcessor) receiveData(data []byte) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.PingPacketMagic:
		log.Debug("initial ping received, sending ping back")
		mp.writeToUsb(packet.NewPingPacketAsBytes())
		return
	case packet.SyncPacketMagic:
		mp.handleSyncPacket(data[4:])
		return
	case packet.AsynPacketMagic:
		mp.handleAsyncPacket(data[4:])
		return
	default:
		log.Warnf("received unknown packet type: %x", data[:4])
	}

	var stop interface{}
	mp.stopSignal <- stop
}

func (mp *messageProcessor) handleSyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data[8:]) {
	case packet.CWPA:
		log.Debug("Received Sync CWPA")
		deviceInfo := packet.NewAsynHpd1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debug("Sending ASYN HPD1")
		mp.writeToUsb(deviceInfo)

		deviceInfo1 := packet.NewAsynHpa1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debug("Sending ASYN HPA1")
		mp.writeToUsb(deviceInfo1)
	case packet.CVRP:
		log.Debug("Received Sync CVRP")
		payload, err := packet.ExtractDictFromBytes(data)
		if err != nil {
			log.Error("Error parsing CVRP packet", err)
			return
		}
		log.Debugf("CVRP:%s", payload.Payload.String())
	default:
		log.Warnf("received unknown sync packet type: %x", data)
	}
}

func (mp *messageProcessor) handleAsyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.FEED:
		log.Debug("Sending FEED")
		mp.writeToUsb(packet.AsynNeedPacketBytes)
	default:
		log.Warnf("received unknown async packet type: %x", data)
	}
}
