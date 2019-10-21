package coremedia

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioStreamBasicDescriptionSerializer(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/adsb-from-hpa-dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, 56)
	adsb := AudioStreamBasicDescription{FormatFlags: 12,
		BytesPerPacket: 4, FramesPerPacket: 1, BytesPerFrame: 4, ChannelsPerFrame: 2, BitsPerChannel: 16, Reserved: 0,
		SampleRate: 48000}
	adsb.SerializeAudioStreamBasicDescription(buffer)

	assert.Equal(t, dat, buffer)

	parsedAdsb, err := NewAudioStreamBasicDescriptionFromBytes(buffer)
	assert.Equal(t, adsb.String(), parsedAdsb.String())
}
