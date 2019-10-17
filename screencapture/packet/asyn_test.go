package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestAsynNeed(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-need")
	if err != nil {
		log.Fatal(err)
	}

	needBytes := packet.AsynNeedPacketBytes(0x0000000102c16ca0)
	assert.Equal(t, dat, needBytes)
}
