package dict_test

import (
	"encoding/hex"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

const ppsHex = "27640033AC5680470133E69E6E04040404"
const spsHex = "28EE3CB0"

func TestParseFormatDescriptor(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/formatdescriptor.bin")
	if err != nil {
		log.Fatal(err)
	}
	fdsc, err := dict.NewFormatDescriptorFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, dict.MediaTypeVideo, fdsc.MediaType)
		assert.Equal(t, decodeSafe(ppsHex), fdsc.PPS)
		assert.Equal(t, decodeSafe(spsHex), fdsc.SPS)
		print(fdsc.String())
	}
}

func decodeSafe(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
