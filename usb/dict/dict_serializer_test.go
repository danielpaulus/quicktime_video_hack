package dict

import (
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
	stringKeyDict := StringKeyDict{Entries: make([]StringKeyEntry, 1)}
	stringKeyDict.Entries[0] = StringKeyEntry{
		Key:   "Valeria",
		Value: true,
	}
	serializedDict := SerializeStringKeyDict(stringKeyDict)
	assert.Equal(t, dat, serializedDict)
}
