package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//AsynFeedPacket contains a CMSampleBuffer and there the actual video data
type AsynFeedPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	CMSampleBuf coremedia.CMSampleBuffer
}

//AsynEatPacket contains a CMSampleBuffer with audio data
type AsynEatPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	CMSampleBuf coremedia.CMSampleBuffer
}

//NewAsynEatPacketFromBytes parses a new AsynEatPacket from bytes
func NewAsynEatPacketFromBytes(data []byte) (AsynEatPacket, error) {
	clockRef, sBuf, err := newAsynCmSampleBufferPacketFromBytes(data, EAT)
	if err != nil {
		return AsynEatPacket{}, err
	}
	return AsynEatPacket{AsyncMagic: AsynPacketMagic, ClockRef: clockRef, MessageType: EAT, CMSampleBuf: sBuf}, nil
}

//NewAsynFeedPacketFromBytes parses a new AsynFeedPacket from bytes
func NewAsynFeedPacketFromBytes(data []byte) (AsynFeedPacket, error) {
	clockRef, sBuf, err := newAsynCmSampleBufferPacketFromBytes(data, FEED)
	if err != nil {
		return AsynFeedPacket{}, err
	}
	return AsynFeedPacket{AsyncMagic: AsynPacketMagic, ClockRef: clockRef, MessageType: FEED, CMSampleBuf: sBuf}, nil
}

func newAsynCmSampleBufferPacketFromBytes(data []byte, magic uint32) (CFTypeID, coremedia.CMSampleBuffer, error) {

	asyncMagic := binary.LittleEndian.Uint32(data)
	if asyncMagic != AsynPacketMagic {
		return 0, coremedia.CMSampleBuffer{}, fmt.Errorf("invalid asyn magic: %x", data)
	}
	clockRef := binary.LittleEndian.Uint64(data[4:])
	messageType := binary.LittleEndian.Uint32(data[12:])
	if messageType != magic {
		return 0, coremedia.CMSampleBuffer{}, fmt.Errorf("invalid packet type in asyn cmsamplebufferpacket:%x", data)
	}
	var cMSampleBuf coremedia.CMSampleBuffer
	var err error
	if magic == FEED {
		cMSampleBuf, err = coremedia.NewCMSampleBufferFromBytesVideo(data[16:])
		if err != nil {
			return 0, coremedia.CMSampleBuffer{}, err
		}
	} else {
		cMSampleBuf, err = coremedia.NewCMSampleBufferFromBytesAudio(data[16:])
		if err != nil {
			return 0, coremedia.CMSampleBuffer{}, err
		}
	}

	return clockRef, cMSampleBuf, nil
}

func (sp AsynFeedPacket) String() string {
	return fmt.Sprintf("ASYN_FEED{ClockRef:%x, sBuf:%s}", sp.ClockRef, sp.CMSampleBuf.String())
}

func (sp AsynEatPacket) String() string {
	return fmt.Sprintf("ASYN_EAT!{ClockRef:%x, sBuf:%s}", sp.ClockRef, sp.CMSampleBuf.String())
}
