package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestTbas(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-tbas")
	if err != nil {
		log.Fatal(err)
	}
	tbasPacket, err := packet.NewAsynTbasPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), tbasPacket.ClockRef)
		assert.Equal(t, uint64(0x1024490c0), tbasPacket.SomeOtherRef)
		assert.Equal(t, "ASYN_TBAS{ClockRef:11123bc18, UnknownRef:1024490c0}", tbasPacket.String())
	}
}
