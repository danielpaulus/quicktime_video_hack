package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
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
	testSerializationOfTimeReply(timePacket, t)
}

func testSerializationOfTimeReply(timePacket packet.SyncTimePacket, t *testing.T) {
	cmtime := coremedia.CMTime{
		CMTimeValue: 0x0000BA62C442E1E1,
		CMTimeScale: 0x3B9ACA00,
		CMTimeFlags: coremedia.KCMTimeFlagsHasBeenRounded,
		CMTimeEpoch: 0,
	}
	replyBytes, err := timePacket.NewReply(cmtime)
	if assert.NoError(t, err) {
		expectedReplyBytes, err := ioutil.ReadFile("fixtures/time-reply1")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, expectedReplyBytes, replyBytes)
	}
}
