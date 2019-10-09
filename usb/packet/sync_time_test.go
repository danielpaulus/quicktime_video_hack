package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/time-request1")
	if err != nil {
		log.Fatal(err)
	}
	timePacket, err := packet.NewSyncTimePacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa67cc17980), timePacket.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, timePacket.SyncMagic)
		assert.Equal(t, packet.TIME, timePacket.MessageType)
		assert.Equal(t, uint64(0x113223d50), timePacket.CorrelationID)
	}
	//testSerializationOfTimeReply(clok, t)
}

func testSerializationOfTimeReply(clok packet.SyncClokPacket, t *testing.T) {
	var clockRef packet.CFTypeID = 0x00007FA67CC17980
	replyBytes := clok.NewReply(clockRef)
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/clok-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
