package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestFeed(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}
	sprpPacket, err := packet.NewAsynFeedPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7ffb5cc32f60), sprpPacket.ClockRef)
		assert.Equal(t, packet.AsynPacketMagic, sprpPacket.AsyncMagic)
		assert.Equal(t, packet.FEED, sprpPacket.MessageType)
	}
}
