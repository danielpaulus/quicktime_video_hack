package usb

import (
	"encoding/binary"
	"encoding/hex"
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
	log.Debugf("Rcv:\n%s", hex.Dump(data))
	//TODO: extractFrame(data)
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
	switch binary.LittleEndian.Uint32(data) {
	case packet.CWPA:
		deviceInfo := packet.NewAsynHpd1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debugf("sending: %s", hex.Dump(deviceInfo))
		mp.writeToUsb(deviceInfo)

		deviceInfo1 := packet.NewAsynHpa1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debugf("sending: %s", hex.Dump(deviceInfo1))
		mp.writeToUsb(deviceInfo1)
	case packet.CVRP:
	default:
		log.Warnf("received unknown sync packet type: %x", data[:4])
	}
}

func (mp *messageProcessor) handleAsyncPacket(data []byte) {
	switch binary.LittleEndian.Uint32(data) {
	case packet.FEED:
		log.Debugf("sending: %s", hex.Dump(packet.AsynNeedPacketBytes))
		mp.writeToUsb(packet.AsynNeedPacketBytes)
	default:
		log.Warnf("received unknown async packet type: %x", data[:4])
	}
}
