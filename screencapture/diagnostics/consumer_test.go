package diagnostics_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/diagnostics"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	waiter := WriteWaiter{make(chan []byte, 100)}

	d := diagnostics.NewDiagnosticsConsumer(waiter, time.Microsecond)
	header := <-waiter.written
	assert.Equal(t, diagnostics.CSVHeader, string(header))
	audioBytes := 35
	videoBytes := 89
	audiobuf := coremedia.CMSampleBuffer{MediaType: coremedia.MediaTypeSound, SampleData: make([]byte, audioBytes)}
	videobuf := coremedia.CMSampleBuffer{MediaType: coremedia.MediaTypeVideo, SampleData: make([]byte, videoBytes)}
	d.Consume(audiobuf)
	d.Consume(videobuf)
	data := <-waiter.written
	d.Stop()
	result := strings.Split(string(data), ",")
	assert.Equal(t, result[0], "1")
	assert.Equal(t, result[1], strconv.Itoa(audioBytes))
	assert.Equal(t, result[2], "1")
	assert.Equal(t, result[3], strconv.Itoa(videoBytes))
}

type WriteWaiter struct {
	written chan []byte
}

func (w WriteWaiter) Write(p []byte) (n int, err error) {
	w.written <- p
	return len(p), nil
}
