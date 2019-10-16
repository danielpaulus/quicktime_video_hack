package dict_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/dict"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/messages"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestBooleanSerialization(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/bulvalue.bin")
	if err != nil {
		log.Fatal(err)
	}
	stringKeyDict := dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 1)}
	stringKeyDict.Entries[0] = dict.StringKeyEntry{
		Key:   "Valeria",
		Value: true,
	}
	serializedDict := dict.SerializeStringKeyDict(stringKeyDict)
	assert.Equal(t, dat, serializedDict)
}

func TestFullSerialization(t *testing.T) {
	dictBytes, err := ioutil.ReadFile("fixtures/serialize_dict.bin")
	if err != nil {
		log.Fatal(err)
	}

	serializedBytes := dict.SerializeStringKeyDict(messages.CreateHpa1DeviceInfoDict())
	assert.Equal(t, dictBytes, serializedBytes)

	dictBytes2, err := ioutil.ReadFile("fixtures/dict.bin")
	if err != nil {
		log.Fatal(err)
	}

	serializedBytes2 := dict.SerializeStringKeyDict(messages.CreateHpd1DeviceInfoDict())
	assert.Equal(t, dictBytes2, serializedBytes2)

}
