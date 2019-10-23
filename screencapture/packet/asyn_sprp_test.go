package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestSprp(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-sprp")
	if err != nil {
		log.Fatal(err)
	}
	sprpPacket, err := packet.NewAsynSprpPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), sprpPacket.ClockRef)
		assert.Equal(t, "ObeyEmptyMediaMarkers", sprpPacket.Property.Key)
		assert.Equal(t, true, sprpPacket.Property.Value.(bool))
		assert.Equal(t, "ASYN_SPRP{ClockRef:11123bc18, Property:{ObeyEmptyMediaMarkers:%!s(bool=true)}}", sprpPacket.String())
	}
}
