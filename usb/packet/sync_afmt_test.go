package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestAfmt(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/afmt-request")
	if err != nil {
		log.Fatal(err)
	}
	afmtPacket, err := packet.NewSyncAfmtPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa66ce20cb0), afmtPacket.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, afmtPacket.SyncMagic)
		assert.Equal(t, packet.AFMT, afmtPacket.MessageType)
		assert.Equal(t, messages.LpcmMagic, afmtPacket.LpcmMagic)
		assert.Equal(t, uint32(0x4c), afmtPacket.LpcmData.Unknown_int1)
		testSerializationOfAfmtReply(afmtPacket, t)
	}
}

func testSerializationOfAfmtReply(clok packet.SyncAfmtPacket, t *testing.T) {
	replyBytes := clok.NewReply()
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/afmt-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
