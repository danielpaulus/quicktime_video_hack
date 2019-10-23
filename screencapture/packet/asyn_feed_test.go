package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}
	feedPacket, err := packet.NewAsynFeedPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7ffb5cc32f60), feedPacket.ClockRef)
	}
}

func TestEat(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-eat")
	if err != nil {
		log.Fatal(err)
	}
	feedPacket, err := packet.NewAsynEatPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x133959728), feedPacket.ClockRef)
	}
}
