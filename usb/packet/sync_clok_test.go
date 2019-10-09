package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestClok(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/clok-request")
	if err != nil {
		log.Fatal(err)
	}
	clok, err := packet.NewSyncClokPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa66cd10250), clok.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, clok.SyncMagic)
		assert.Equal(t, packet.CLOK, clok.MessageType)
		assert.Equal(t, uint64(0x113584970), clok.CorrelationID)
	}
	testSerializationOfClokReply(clok, t)
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
