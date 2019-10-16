package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestSprp(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-sprp")
	if err != nil {
		log.Fatal(err)
	}
	sprpPacket, err := packet.NewAsynSprpPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), sprpPacket.ClockRef)
		assert.Equal(t, packet.AsynPacketMagic, sprpPacket.AsyncMagic)
		assert.Equal(t, packet.SPRP, sprpPacket.MessageType)
		assert.Equal(t, "ObeyEmptyMediaMarkers", sprpPacket.Property.Key)
		assert.Equal(t, true, sprpPacket.Property.Value.(bool))

	}
}
