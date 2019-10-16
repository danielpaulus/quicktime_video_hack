package messages

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestLpcmSerializer(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/lpcm.bin")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, dat, createLpcmInfo())

}
