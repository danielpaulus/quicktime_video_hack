package packet

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestSimpleDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/sync_cvrp.bin")
	if err != nil {
		log.Fatal(err)
	}
	syn, err := ExtractDictFromBytes(dat[4:])
	if assert.NoError(t, err) {
		const expectedHeader uint64 = 1
		assert.Equal(t, expectedHeader, syn.Header)
		assert.Equal(t, CVRP, syn.Magic)
	}
}
