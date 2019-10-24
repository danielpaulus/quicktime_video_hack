package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestCwpa(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/cwpa-request1")
	if err != nil {
		log.Fatal(err)
	}
	cwpa, err := packet.NewSyncCwpaPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, packet.EmptyCFType, cwpa.ClockRef)
		assert.Equal(t, uint64(0x1135a74e0), cwpa.DeviceClockRef)
		assert.Equal(t, uint64(0x113573de0), cwpa.CorrelationID)
		assert.Equal(t, "SYNC_CWPA{ClockRef:1, CorrelationID:113573de0, DeviceClockRef:1135a74e0}", cwpa.String())
	}
	_, err = packet.NewSyncCwpaPacketFromBytes(dat)
	assert.Error(t, err)

	brokenMessage := make([]byte, len(dat))
	copy(brokenMessage, dat)
	testIncorrectClockRefProducesError(brokenMessage, t)

	copy(brokenMessage, dat)
	testIncorrectSubtypeProducesError(brokenMessage, t, cwpa)
	testSerializationOfReply(cwpa, t)
}

func testIncorrectClockRefProducesError(brokenMessage []byte, t *testing.T) {
	brokenMessage[9] = 0xFF
	_, err := packet.NewSyncCwpaPacketFromBytes(brokenMessage[4:])
	assert.Error(t, err)
}

func testIncorrectSubtypeProducesError(brokenMessage []byte, t *testing.T, cwpa packet.SyncCwpaPacket) {
	brokenMessage[17] = 0xFF
	_, err := packet.NewSyncCwpaPacketFromBytes(brokenMessage[4:])
	assert.Error(t, err)
}

func testSerializationOfReply(cwpa packet.SyncCwpaPacket, t *testing.T) {
	var clockRef packet.CFTypeID = 0x00007FA66CE20CB0
	replyBytes := cwpa.NewReply(clockRef)
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/cwpa-reply1")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
