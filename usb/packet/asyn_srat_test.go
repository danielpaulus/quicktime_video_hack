package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestSrat(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-srat")
	if err != nil {
		log.Fatal(err)
	}
	sratPacket, err := packet.NewAsynSratPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), sratPacket.ClockRef)
		assert.Equal(t, packet.AsynPacketMagic, sratPacket.AsyncMagic)
		assert.Equal(t, packet.SRAT, sratPacket.MessageType)
		assert.Equal(t, float32(1), sratPacket.Rate1)
		assert.Equal(t, float32(1), sratPacket.Rate2)
		assert.Equal(t, uint32(1000000000), sratPacket.Time.CMTimeScale)

	}
}
