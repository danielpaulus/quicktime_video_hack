package coremedia_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

func TestWav(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/out.raw")
	if err != nil {
		log.Fatal(err)
	}

	headerPlaceholder := make([]byte, 44)

	file, err := os.Create("/Users/danielpaulus/tmp/out1.wav")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(headerPlaceholder)
	defer file.Close()
	file.Write(dat)
	buffer := bytes.NewBuffer(make([]byte, 100))
	buffer.Reset()

	riffHeader := coremedia.NewRiffHeader(len(dat))
	riffHeader.Serialize(buffer)

	fmtSubChunk := coremedia.NewFmtSubChunk()
	fmtSubChunk.Serialize(buffer)

	coremedia.WriteWavDataSubChunkHeader(buffer, len(dat))
	file.WriteAt(buffer.Bytes(), 0)
	//log.Fatal(fmt.Errorf("test:%x", buffer.Bytes()))
}
