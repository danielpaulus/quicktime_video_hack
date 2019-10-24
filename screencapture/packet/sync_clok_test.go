package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestClok(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/clok-request")
	if err != nil {
		log.Fatal(err)
	}
	clok, err := packet.NewSyncClokPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa66cd10250), clok.ClockRef)
		assert.Equal(t, uint64(0x113584970), clok.CorrelationID)
		assert.Equal(t, "SYNC_CLOK{ClockRef:7fa66cd10250, CorrelationID:113584970}", clok.String())
	}
	testSerializationOfClokReply(clok, t)
	_, err = packet.NewSyncClokPacketFromBytes(dat)
	assert.Error(t, err)
}

func testSerializationOfClokReply(clok packet.SyncClokPacket, t *testing.T) {
	var clockRef packet.CFTypeID = 0x00007FA67CC17980
	replyBytes := clok.NewReply(clockRef)
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/clok-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
