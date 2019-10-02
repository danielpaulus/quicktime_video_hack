package usb

import (
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

var i int

func (mp *messageProcessor) receiveData(data []byte) {
	log.Debugf("Rcv:\n%s", hex.Dump(data))
	//TODO: extractFrame(data)
	if mp.connectionState == initialState {
		log.Debug("initial ping received, sending ping back")
		mp.respondToPing(packet.NewPingPacketAsBytes())

		mp.connectionState = pingSent
		return
	}

	if i == 0 {
		deviceInfo := packet.NewAsynHpd1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debugf("sending: %s", hex.Dump(deviceInfo))
		mp.writeToUsb(deviceInfo)
		i++
		return
	}
	if i == 1 {
		deviceInfo := packet.NewAsynHpa1Packet(messages.CreateHpd1DeviceInfoDict())
		log.Debugf("sending: %s", hex.Dump(deviceInfo))
		mp.writeToUsb(deviceInfo)
		i++
		return
	}
	if i == 2 {
		log.Debugf("sending: %s", hex.Dump(packet.AsynNeedPacketBytes))
		mp.writeToUsb(packet.AsynNeedPacketBytes)
		return
	}

	//deviceInfo2 := packet.NewAsynHpa1Packet(messages.CreateHpa1DeviceInfoDict())

	var stop interface{}
	mp.stopSignal <- stop
}

func (mp messageProcessor) respondToPing(bytes []byte) {
	mp.writeToUsb(bytes)
}
