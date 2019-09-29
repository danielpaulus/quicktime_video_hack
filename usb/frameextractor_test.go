package usb_test

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/usb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompletePacket(t *testing.T) {
	small := make([]byte, 30)
	binary.LittleEndian.PutUint32(small, 30)
	small[29] = 5
	fe := usb.NewLengthFieldBasedFrameExtractor()
	frame, frameReturned := fe.ExtractFrame(small)
	assert.True(t, frameReturned)
	assert.Equal(t, small[4:], frame)
	frame1, frameReturned2 := fe.ExtractFrame(small)
	assert.True(t, frameReturned2)
	assert.Equal(t, small[4:], frame1)
}

func TestIncompletePacket(t *testing.T) {
	data := make([]byte, 30)
	binary.LittleEndian.PutUint32(data, 40)
	data[29] = 5
	remaining := make([]byte, 10)
	fe := usb.NewLengthFieldBasedFrameExtractor()
	_, frameReturned := fe.ExtractFrame(data)
	assert.False(t, frameReturned)

	frame, frameReturned2 := fe.ExtractFrame(remaining)
	assert.True(t, frameReturned2)
	assert.Equal(t, append(data, remaining...)[4:], frame)

}
