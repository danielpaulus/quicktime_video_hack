package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestAfmt(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/afmt-request")
	if err != nil {
		log.Fatal(err)
	}
	afmtPacket, err := packet.NewSyncAfmtPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x11123bc18), afmtPacket.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, afmtPacket.SyncMagic)
		assert.Equal(t, packet.AFMT, afmtPacket.MessageType)
		assert.Equal(t, messages.LpcmMagic, afmtPacket.LpcmMagic)

	}
}
