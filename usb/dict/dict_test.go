package dict_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestBooleanEntry(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/bulvalue.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := dict.NewDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(mydict.Entries))
		assert.Equal(t, "Valeria", mydict.Entries[0].Key)
		assert.Equal(t, true, mydict.Entries[0].Value.(bool))
	}
}

func TestSimpleDictEntry(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := dict.NewDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(mydict.Entries))
		assert.Equal(t, "Valeria", mydict.Entries[0].Key)
		assert.Equal(t, true, mydict.Entries[0].Value.(bool))

		assert.Equal(t, "HEVCDecoderSupports444", mydict.Entries[1].Key)
		assert.Equal(t, true, mydict.Entries[1].Value.(bool))

		assert.Equal(t, "DisplaySize", mydict.Entries[2].Key)
		displaySize := mydict.Entries[2].Value.(dict.Dict)
		assert.Equal(t, 2, len(displaySize.Entries))

		assert.Equal(t, "Width", displaySize.Entries[0].Key)
		assert.Equal(t, float64(1920), displaySize.Entries[0].Value.(dict.NSNumber).FloatValue)

		assert.Equal(t, "Height", displaySize.Entries[1].Key)
		assert.Equal(t, float64(1200), displaySize.Entries[1].Value.(dict.NSNumber).FloatValue)
	}
}
