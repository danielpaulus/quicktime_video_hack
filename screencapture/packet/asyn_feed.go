package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//AsynCmSampleBufPacket contains a CMSampleBuffer with audio or video data
type AsynCmSampleBufPacket struct {
	ClockRef    CFTypeID
	CMSampleBuf coremedia.CMSampleBuffer
}

//NewAsynCmSampleBufPacketFromBytes parses a new AsynCmSampleBufPacket from bytes
func NewAsynCmSampleBufPacketFromBytes(data []byte) (AsynCmSampleBufPacket, error) {
	clockRef, sBuf, err := newAsynCmSampleBufferPacketFromBytes(data)
	if err != nil {
		return AsynCmSampleBufPacket{}, err
	}
	return AsynCmSampleBufPacket{ClockRef: clockRef, CMSampleBuf: sBuf}, nil
}

func newAsynCmSampleBufferPacketFromBytes(data []byte) (CFTypeID, coremedia.CMSampleBuffer, error) {
	magic := binary.LittleEndian.Uint32(data[12:])
	_, clockRef, err := ParseAsynHeader(data, magic)
	if err != nil {
		return 0, coremedia.CMSampleBuffer{}, err
	}

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

func (sp AsynCmSampleBufPacket) String() string {
	return fmt.Sprintf("ASYN_SBUF{ClockRef:%x, sBuf:%s}", sp.ClockRef, sp.CMSampleBuf.String())
}
