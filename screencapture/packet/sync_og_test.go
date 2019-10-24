package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestOg(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/og-request")
	if err != nil {
		log.Fatal(err)
	}
	og, err := packet.NewSyncOgPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fba35425ff0), og.ClockRef)
		assert.Equal(t, uint64(0x102d32f30), og.CorrelationID)
		assert.Equal(t, "SYNC_OG{ClockRef:7fba35425ff0, CorrelationID:102d32f30, Unknown:1}", og.String())
	}
	testSerializationOfOgReply(og, t)
	_, err = packet.NewSyncOgPacketFromBytes(dat)
	assert.Error(t, err)
}

func testSerializationOfOgReply(clok packet.SyncOgPacket, t *testing.T) {
	replyBytes := clok.NewReply()
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/og-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
