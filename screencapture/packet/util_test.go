package packet_test

import (
	"encoding/binary"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestAsyn(t *testing.T) {
	data, expectedBytes, expectedClockRef := createAsynFeedPacket()

	remainingBytes, clockRef, err := packet.ParseAsynHeader(data, packet.FEED)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedClockRef, clockRef)
		assert.Equal(t, expectedBytes, remainingBytes)
	}

	_, _, err = packet.ParseAsynHeader(data, packet.TJMP)
	assert.Error(t, err)

	//break asyn magic marker
	data[0] = 5
	_, _, err = packet.ParseAsynHeader(data, packet.FEED)
	assert.Error(t, err)
}

func createAsynFeedPacket() ([]byte, []byte, packet.CFTypeID) {
	data := make([]byte, 20)
	binary.LittleEndian.PutUint32(data, packet.AsynPacketMagic)
	expectedClockRef := uint64(0xff1233aabbcc00)
	binary.LittleEndian.PutUint64(data[4:], expectedClockRef)
	binary.LittleEndian.PutUint32(data[12:], packet.FEED)
	expectedBytes := []byte{1, 2, 3, 4}
	copy(data[16:], expectedBytes)
	return data, expectedBytes, expectedClockRef
}
