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
