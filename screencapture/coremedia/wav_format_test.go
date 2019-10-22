package coremedia_test

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

const expectedBytes = "524946461802000057415645666d7420100000000100020080bb000000ee02000400100064617461f401000044616e69656c"

func TestWavHeaderWrittenCorrectly(t *testing.T) {

	headerPlaceholder := make([]byte, 44)

	file, err := ioutil.TempFile("", "golangtemp*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.Write(headerPlaceholder)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(([]byte)("Daniel"))
	if err != nil {
		log.Fatal(err)
	}

	err = coremedia.WriteWavHeader(500, file)
	if assert.NoError(t, err) {

		dat, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Fatal(err)
		}
		expectedBytes, _ := hex.DecodeString(expectedBytes)
		assert.Equal(t, expectedBytes, dat)
	}

}
