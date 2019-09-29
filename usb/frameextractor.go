package usb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

type lengthFieldBasedFrameExtractor struct {
	frameBuffer       *bytes.Buffer
	readyForNextFrame bool
	nextFrameSize     int
}

func NewLengthFieldBasedFrameExtractor() *lengthFieldBasedFrameExtractor {
	extractor := &lengthFieldBasedFrameExtractor{
		frameBuffer:       bytes.NewBuffer(make([]byte, 1024*1024*2)),
		readyForNextFrame: true}
	extractor.frameBuffer.Reset()
	return extractor
}

func (fe *lengthFieldBasedFrameExtractor) ExtractFrame(bytes []byte) ([]byte, bool) {
	if fe.readyForNextFrame && fe.frameBuffer.Len() == 0 {
		return fe.handleNewFrame(bytes)
	}
	if fe.readyForNextFrame && fe.frameBuffer.Len() != 0 {
		fe.nextFrameSize = int(binary.LittleEndian.Uint32(fe.frameBuffer.Next(4))) - 4
		fe.readyForNextFrame = false
		return fe.ExtractFrame(bytes)
	}
	fe.frameBuffer.Write(bytes)
	if fe.frameBuffer.Len() >= fe.nextFrameSize {
		frame := make([]byte, fe.nextFrameSize)
		_, err := fe.frameBuffer.Read(frame)
		if err != nil {
			log.Error("Failed reading from internal buffer", err)
		}
		fe.readyForNextFrame = true
		return frame, true
	}
	return nil, false
}

func (fe *lengthFieldBasedFrameExtractor) handleNewFrame(bytes []byte) ([]byte, bool) {
	if len(bytes) < 4 {
		log.Fatalf("Received less than four bytes, cannot read a valid frameLength field: %s", hex.Dump(bytes))
	}

	frameLength := int(binary.LittleEndian.Uint32(bytes[:4]))
	if len(bytes) == frameLength {
		return bytes[4:], true
	}
	if len(bytes) > frameLength {
		fe.frameBuffer.Write(bytes[frameLength:])
		return bytes[4:frameLength], true
	}
	fe.readyForNextFrame = false
	fe.frameBuffer.Write(bytes[4:])
	fe.nextFrameSize = frameLength - 4
	return nil, false
}
