package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"

	"github.com/stretchr/testify/assert"
)

func TestAsynHP(t *testing.T) {
	hpa0Dat, err := ioutil.ReadFile("fixtures/asyn-hpa0")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, hpa0Dat, packet.NewAsynHPA0(0x0000000102C5FC10))

	hpd0Dat, err := ioutil.ReadFile("fixtures/asyn-hpd0")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, hpd0Dat, packet.NewAsynHPD0())
}
