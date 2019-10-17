package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestOg(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/og-request")
	if err != nil {
		log.Fatal(err)
	}
	clok, err := packet.NewSyncOgPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fba35425ff0), clok.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, clok.SyncMagic)
		assert.Equal(t, packet.OG, clok.MessageType)
		assert.Equal(t, uint64(0x102d32f30), clok.CorrelationID)
	}
	testSerializationOfOgReply(clok, t)
}

func testSerializationOfOgReply(clok packet.SyncOgPacket, t *testing.T) {
	replyBytes := clok.NewReply()
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/og-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
