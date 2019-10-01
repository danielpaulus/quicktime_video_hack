package dict_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestParseFormatDescriptor(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/formatdescriptor.bin")
	if err != nil {
		log.Fatal(err)
	}
	fdsc, err := dict.NewFormatDescriptorFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, dict.MediaTypeVideo, fdsc.MediaType)
		print(fdsc.String())
	}
}
