package packet_test

import (
	"encoding/binary"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestParseAsynHeader(t *testing.T) {
	data, expectedBytes, expectedClockRef := createAsynFeedPacket()

	remainingBytes, clockRef, err := packet.ParseAsynHeader(data, packet.FEED)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedClockRef, clockRef)
		assert.Equal(t, expectedBytes, remainingBytes)
	}

	_, _, err = packet.ParseAsynHeader(data, packet.TJMP)
	assert.Error(t, err)
	assert.Equal(t, "invalid packet type:'deef' - packethex: 6e79736100ccbbaa3312ff006465656601020304", err.Error())
	//break asyn magic marker
	data[0] = 80
	_, _, err = packet.ParseAsynHeader(data, packet.FEED)
	assert.Error(t, err)
	assert.Equal(t, "invalid packet magic 'Pysa' - packethex: 5079736100ccbbaa3312ff006465656601020304", err.Error())
}

func TestParseSyncHeader(t *testing.T) {
	data, expectedBytes, expectedClockRef, expectedCorrelationID := createSyncClokPacket()

	remainingBytes, clockRef, correlationID, err := packet.ParseSyncHeader(data, packet.FEED)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedClockRef, clockRef)
		assert.Equal(t, expectedBytes, remainingBytes)
		assert.Equal(t, expectedCorrelationID, correlationID)
	}

	_, _, err = packet.ParseAsynHeader(data, packet.TJMP)
	assert.Error(t, err)

	//break sync magic marker
	data[0] = 80
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

func createSyncClokPacket() ([]byte, []byte, packet.CFTypeID, uint64) {
	data := make([]byte, 28)
	binary.LittleEndian.PutUint32(data, packet.SyncPacketMagic)
	expectedClockRef := uint64(0xff1233aabbcc00)
	binary.LittleEndian.PutUint64(data[4:], expectedClockRef)
	binary.LittleEndian.PutUint32(data[12:], packet.FEED)
	expectedCorrelationID := uint64(0xFFDDFFAA)
	binary.LittleEndian.PutUint64(data[16:], expectedCorrelationID)
	expectedBytes := []byte{1, 2, 3, 4}
	copy(data[24:], expectedBytes)
	return data, expectedBytes, expectedClockRef, expectedCorrelationID
}
