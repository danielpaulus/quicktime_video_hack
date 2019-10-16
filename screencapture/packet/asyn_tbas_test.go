package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestTbas(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-tbas")
	if err != nil {
		log.Fatal(err)
	}
	sratPacket, err := packet.NewAsynTbasPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), sratPacket.ClockRef)
		assert.Equal(t, packet.AsynPacketMagic, sratPacket.AsyncMagic)
		assert.Equal(t, packet.TBAS, sratPacket.MessageType)
		assert.Equal(t, uint64(0x1024490c0), sratPacket.SomeOtherRef)
	}
}
