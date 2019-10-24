package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestSrat(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-srat")
	if err != nil {
		log.Fatal(err)
	}
	sratPacket, err := packet.NewAsynSratPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), sratPacket.ClockRef)
		assert.Equal(t, float32(1), sratPacket.Rate1)
		assert.Equal(t, float32(1), sratPacket.Rate2)
		assert.Equal(t, uint32(1000000000), sratPacket.Time.CMTimeScale)
		assert.Equal(t, "ASYN_SRAT{ClockRef:11123bc18, Rate1:1.000000, Rate2:1.000000, Time:CMTime{1570648854000190667/1000000000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}}", sratPacket.String())
	}
}
