package screencapture

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

//LengthFieldBasedFrameExtractor extracts frames from the packetized byte stream we get from USB
type LengthFieldBasedFrameExtractor struct {
	frameBuffer       *bytes.Buffer
	readyForNextFrame bool
	nextFrameSize     int
}

//NewLengthFieldBasedFrameExtractor intializes a new Extractor with a 2MB buffer
func NewLengthFieldBasedFrameExtractor() *LengthFieldBasedFrameExtractor {
	extractor := &LengthFieldBasedFrameExtractor{
		frameBuffer:       bytes.NewBuffer(make([]byte, 1024*1024*2)),
		readyForNextFrame: true}
	extractor.frameBuffer.Reset()
	return extractor
}

//ExtractFrame writes new bytes into the extractor and if possible
//returns a frame when the returned bool is true and nil otherwise.
//It can be called with an empty slice to check if there are multiple frames in the Extractor.
func (fe *LengthFieldBasedFrameExtractor) ExtractFrame(bytes []byte) ([]byte, bool) {
	if fe.readyForNextFrame && fe.frameBuffer.Len() == 0 {
		return fe.handleNewFrame(bytes)
	}
	if fe.readyForNextFrame && fe.frameBuffer.Len() != 0 {
		if fe.frameBuffer.Len() < 4 {
			/* examples:
			FATA[0032] wtf:00000000  ac 10                                             |..|
			*/
			/*
				{"level":"fatal","msg":"wtf:00000000  ac                                                |.|\n","time":"2019-11-18T12:09:10+01:00"}
			*/
			/*
				{"level":"fatal","msg":"wtf:00000000  ac 10 00                                          |...|\n","time":"2019-11-18T12:12:11+01:00"}
			*/
			log.Warn("received invalid message, dropping:" + hex.Dump(fe.frameBuffer.Bytes()))
			return nil, false
		}
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

func (fe *LengthFieldBasedFrameExtractor) handleNewFrame(bytes []byte) ([]byte, bool) {
	//ZeroLengthPackages are sometimes sent on USB Connections. It is safe to ignore them.
	if len(bytes) == 0 {
		log.Debugf("skipping 0")
		return nil, false
	}
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
