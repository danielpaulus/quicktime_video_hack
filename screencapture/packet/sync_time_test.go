package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/time-request1")
	if err != nil {
		log.Fatal(err)
	}
	timePacket, err := packet.NewSyncTimePacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa67cc17980), timePacket.ClockRef)
		assert.Equal(t, uint64(0x113223d50), timePacket.CorrelationID)
		assert.Equal(t, "SYNC_TIME{ClockRef:7fa67cc17980, CorrelationID:113223d50}", timePacket.String())
	}
	testSerializationOfTimeReply(timePacket, t)
	_, err = packet.NewSyncTimePacketFromBytes(dat)
	assert.Error(t, err)
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
