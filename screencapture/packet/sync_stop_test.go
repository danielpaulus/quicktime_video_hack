package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestStop(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/stop-request")
	if err != nil {
		log.Fatal(err)
	}
	stop, err := packet.NewSyncStopPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fba35425ff0), stop.ClockRef)
		assert.Equal(t, uint64(0x102fd4910), stop.CorrelationID)
		assert.Equal(t, "SYNC_STOP{ClockRef:7fba35425ff0, CorrelationID:102fd4910}", stop.String())
	}
	testSerializationOfStopReply(stop, t)
	_, err = packet.NewSyncStopPacketFromBytes(dat)
	assert.Error(t, err)

}

func testSerializationOfStopReply(stop packet.SyncStopPacket, t *testing.T) {
	replyBytes := stop.NewReply()
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/stop-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
