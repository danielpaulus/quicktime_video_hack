package diagnostics_test

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/diagnostics"
)

func TestConsumer(t *testing.T) {
	buf := new(bytes.Buffer)
	d := diagnostics.NewDiagnosticsConsumer(buf, time.Microsecond)
	audioBytes := 35
	videoBytes := 89
	audiobuf := coremedia.CMSampleBuffer{MediaType: coremedia.MediaTypeSound, SampleData: make([]byte, audioBytes)}
	videobuf := coremedia.CMSampleBuffer{MediaType: coremedia.MediaTypeVideo, SampleData: make([]byte, videoBytes)}
	d.Consume(audiobuf)
	d.Consume(videobuf)
	time.Sleep(time.Second)
	d.Stop()
	result := strings.Split(string(buf.Bytes()), "\n")
	log.Fatal(result)
}
