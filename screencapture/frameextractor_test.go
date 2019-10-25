package screencapture_test

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompletePacket(t *testing.T) {
	small := make([]byte, 30)
	binary.LittleEndian.PutUint32(small, 30)
	small[29] = 5
	fe := screencapture.NewLengthFieldBasedFrameExtractor()
	frame, frameReturned := fe.ExtractFrame(small)
	assert.True(t, frameReturned)
	assert.Equal(t, small[4:], frame)
	frame1, frameReturned2 := fe.ExtractFrame(small)
	assert.True(t, frameReturned2)
	assert.Equal(t, small[4:], frame1)

}

func TestZeroLengthPacket(t *testing.T) {
	fe := screencapture.NewLengthFieldBasedFrameExtractor()
	_, frameReturned := fe.ExtractFrame(make([]byte, 0))
	assert.False(t, frameReturned)
}

func TestIncompletePacket(t *testing.T) {
	data := make([]byte, 30)
	binary.LittleEndian.PutUint32(data, 40)
	data[29] = 5
	remaining := make([]byte, 10)
	fe := screencapture.NewLengthFieldBasedFrameExtractor()
	_, frameReturned := fe.ExtractFrame(data)
	assert.False(t, frameReturned)

	frame, frameReturned2 := fe.ExtractFrame(remaining)
	assert.True(t, frameReturned2)
	assert.Equal(t, append(data, remaining...)[4:], frame)
}

func TestIncompletePacketWithFollowUp(t *testing.T) {
	data := make([]byte, 30)
	binary.LittleEndian.PutUint32(data, 40)
	data[29] = 5
	remaining := make([]byte, 10)

	small := make([]byte, 30)
	binary.LittleEndian.PutUint32(small, 30)
	small[29] = 5

	fe := screencapture.NewLengthFieldBasedFrameExtractor()
	_, frameReturned := fe.ExtractFrame(data)
	assert.False(t, frameReturned)

	frame, frameReturned2 := fe.ExtractFrame(append(remaining, small...))
	assert.True(t, frameReturned2)
	assert.Equal(t, append(data, remaining...)[4:], frame)

	frame2, frameReturned3 := fe.ExtractFrame(make([]byte, 0))
	assert.True(t, frameReturned3)
	assert.Equal(t, small[4:], frame2)
}

func TestMultipleCompletePacketsWithIncompleteFollowUp(t *testing.T) {
	firstFrame := make([]byte, 30)
	binary.LittleEndian.PutUint32(firstFrame, 30)
	firstFrame[29] = 5

	secondFrame := make([]byte, 37)
	binary.LittleEndian.PutUint32(secondFrame, 37)
	firstFrame[29] = 7

	thirdFrame := make([]byte, 57)
	binary.LittleEndian.PutUint32(thirdFrame, 57)
	thirdFrame[29] = 37

	remaining := make([]byte, 10)
	binary.LittleEndian.PutUint32(remaining, 570)

	fe := screencapture.NewLengthFieldBasedFrameExtractor()
	frame, frameReturned := fe.ExtractFrame(append(firstFrame, secondFrame...))
	assert.True(t, frameReturned)
	assert.Equal(t, firstFrame[4:], frame)

	frame2, _ := fe.ExtractFrame(append(thirdFrame, remaining...))
	assert.Equal(t, secondFrame[4:], frame2)

	frame3, _ := fe.ExtractFrame(make([]byte, 0))
	assert.Equal(t, thirdFrame[4:], frame3)

	_, frameReturned2 := fe.ExtractFrame(make([]byte, 0))
	assert.False(t, frameReturned2)
}
