package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestAfmt(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/afmt-request")
	if err != nil {
		log.Fatal(err)
	}
	afmtPacket, err := packet.NewSyncAfmtPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fa66ce20cb0), afmtPacket.ClockRef)
		expectedAsbd := coremedia.DefaultAudioStreamBasicDescription()
		expectedAsbd.FormatFlags = 0x4C
		assert.Equal(t, expectedAsbd, afmtPacket.AudioStreamBasicDescription)
		expectedString := "SYNC_AFMT{ClockRef:7fa66ce20cb0, CorrelationID:113229d80, AudioStreamBasicDescription:{SampleRate:48000.000000,FormatFlags:76,BytesPerPacket:4,FramesPerPacket:1,BytesPerFrame:4,ChannelsPerFrame:2,BitsPerChannel:16,Reserved:0}}"
		assert.Equal(t, expectedString, afmtPacket.String())
		testSerializationOfAfmtReply(afmtPacket, t)
	}

	_, err = packet.NewSyncAfmtPacketFromBytes(dat)
	assert.Error(t, err)
}

func testSerializationOfAfmtReply(clok packet.SyncAfmtPacket, t *testing.T) {
	replyBytes := clok.NewReply()
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/afmt-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
