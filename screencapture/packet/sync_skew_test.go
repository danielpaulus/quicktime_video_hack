package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestSkew(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/skew-request")
	if err != nil {
		log.Fatal(err)
	}
	skew, err := packet.NewSyncSkewPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fba35425ff0), skew.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, skew.SyncMagic)
		assert.Equal(t, packet.SKEW, skew.MessageType)
		assert.Equal(t, uint64(0x102fdb960), skew.CorrelationID)
		assert.Equal(t, "SYNC_SKEW{ClockRef:7fba35425ff0, CorrelationID:102fdb960}", skew.String())
	}
	testSerializationOfSkewReply(skew, t)
	
}

func testSerializationOfSkewReply(skew packet.SyncSkewPacket, t *testing.T) {
	replyBytes := skew.NewReply(float64(48000))
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/skew-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}