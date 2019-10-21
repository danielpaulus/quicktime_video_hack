package coremedia_test

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

const ppsHex = "27640033AC5680470133E69E6E04040404"
const spsHex = "28EE3CB0"

func TestParseFormatDescriptor(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/formatdescriptor.bin")
	if err != nil {
		log.Fatal(err)
	}
	fdsc, err := coremedia.NewFormatDescriptorFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, coremedia.MediaTypeVideo, fdsc.MediaType)
		assert.Equal(t, decodeSafe(ppsHex), fdsc.PPS)
		assert.Equal(t, decodeSafe(spsHex), fdsc.SPS)
		println(fdsc.String())
	}
}

func decodeSafe(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func TestParseFormatDescriptorAudio(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/formatdescriptor-audio.bin")
	if err != nil {
		log.Fatal(err)
	}
	fdsc, err := coremedia.NewFormatDescriptorFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, coremedia.MediaTypeSound, fdsc.MediaType)
		println(fdsc.String())
	}
}
