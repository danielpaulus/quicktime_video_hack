package coremedia_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

func TestBooleanSerialization(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/bulvalue.bin")
	if err != nil {
		log.Fatal(err)
	}
	stringKeyDict := coremedia.StringKeyDict{Entries: make([]coremedia.StringKeyEntry, 1)}
	stringKeyDict.Entries[0] = coremedia.StringKeyEntry{
		Key:   "Valeria",
		Value: true,
	}
	serializedDict := coremedia.SerializeStringKeyDict(stringKeyDict)
	assert.Equal(t, dat, serializedDict)
}

func TestFullSerialization(t *testing.T) {
	dictBytes, err := ioutil.ReadFile("fixtures/serialize_dict.bin")
	if err != nil {
		log.Fatal(err)
	}

	serializedBytes := coremedia.SerializeStringKeyDict(packet.CreateHpa1DeviceInfoDict())
	assert.Equal(t, dictBytes, serializedBytes)

	dictBytes2, err := ioutil.ReadFile("fixtures/dict.bin")
	if err != nil {
		log.Fatal(err)
	}

	serializedBytes2 := coremedia.SerializeStringKeyDict(packet.CreateHpd1DeviceInfoDict())
	assert.Equal(t, dictBytes2, serializedBytes2)

}
