package usb

import (
	"encoding/hex"
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
	if mp.connectionState == initialState {
		log.Debug("initial ping received, sending ping back")
		mp.respondToPing(data)
		mp.connectionState = pingSent
		return
	}

	var stop interface{}
	mp.stopSignal <- stop
}

func (mp messageProcessor) respondToPing(bytes []byte) {
	mp.writeToUsb(bytes)
}
