package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestRels(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-rels")
	if err != nil {
		log.Fatal(err)
	}
	relsPacket, err := packet.NewAsynRelsPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7fba35608a00), relsPacket.ClockRef)
		assert.Equal(t, "ASYN_RELS{ClockRef:7fba35608a00}", relsPacket.String())
	}
}
