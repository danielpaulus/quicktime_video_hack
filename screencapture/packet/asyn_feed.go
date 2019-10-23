package packet

import (
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//AsynFeedPacket contains a CMSampleBuffer and there the actual video data
type AsynFeedPacket struct {
	ClockRef    CFTypeID
	CMSampleBuf coremedia.CMSampleBuffer
}

//AsynEatPacket contains a CMSampleBuffer with audio data
type AsynEatPacket struct {
	ClockRef    CFTypeID
	CMSampleBuf coremedia.CMSampleBuffer
}

//NewAsynEatPacketFromBytes parses a new AsynEatPacket from bytes
func NewAsynEatPacketFromBytes(data []byte) (AsynEatPacket, error) {
	clockRef, sBuf, err := newAsynCmSampleBufferPacketFromBytes(data, EAT)
	if err != nil {
		return AsynEatPacket{}, err
	}
	return AsynEatPacket{ClockRef: clockRef, CMSampleBuf: sBuf}, nil
}

//NewAsynFeedPacketFromBytes parses a new AsynFeedPacket from bytes
func NewAsynFeedPacketFromBytes(data []byte) (AsynFeedPacket, error) {
	clockRef, sBuf, err := newAsynCmSampleBufferPacketFromBytes(data, FEED)
	if err != nil {
		return AsynFeedPacket{}, err
	}
	return AsynFeedPacket{ClockRef: clockRef, CMSampleBuf: sBuf}, nil
}

func newAsynCmSampleBufferPacketFromBytes(data []byte, magic uint32) (CFTypeID, coremedia.CMSampleBuffer, error) {
	_, clockRef, err := ParseAsynHeader(data, magic)
	if err != nil {
		return 0, coremedia.CMSampleBuffer{}, err
	}
	clockRef = clockRef

	var cMSampleBuf coremedia.CMSampleBuffer

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
