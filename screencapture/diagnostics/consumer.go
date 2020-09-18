package diagnostics

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

type DiagnosticsConsumer struct {
	file            io.Writer
	audioSamplesRcv uint64
	videoSamplesRcv uint64
	audioBytesRcv   uint64
	videoBytesRcv   uint64
	mux             sync.Mutex
}

func NewDiagnosticsConsumer(outfile io.Writer) *DiagnosticsConsumer {
	d := &DiagnosticsConsumer{file: outfile}
	go fileWriter(d)
	return d
}

func fileWriter(d *DiagnosticsConsumer) {
	d.file.Write([]byte("audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv, heapobjects, alloc"))

	for {
		time.Sleep(time.Second * 10)
		audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv := readAndReset(d)
		heapobjects, alloc := getMemStats()
		csvLine := fmt.Sprintf("%d,%d,%d,%d,%d,%d", audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv, heapobjects, alloc)
		_, err := d.file.Write([]byte(csvLine))
		if err != nil {
			log.Fatalf("Failed writing to metricsfile:%+v", err)
		}
	}
}

func getMemStats() (uint64, uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapObjects, m.Alloc
}

func readAndReset(d *DiagnosticsConsumer) (uint64, uint64, uint64, uint64) {
	d.mux.Lock()
	defer d.mux.Unlock()
	audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv := d.audioSamplesRcv, d.audioBytesRcv, d.videoSamplesRcv, d.videoBytesRcv
	d.audioSamplesRcv, d.audioBytesRcv, d.videoSamplesRcv, d.videoBytesRcv = 0, 0, 0, 0
	return audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv
}

func (d *DiagnosticsConsumer) Consume(buf coremedia.CMSampleBuffer) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	if buf.MediaType == coremedia.MediaTypeSound {
		return d.consumeAudio(buf)
	}

	return d.consumeVideo(buf)
}

func (d *DiagnosticsConsumer) consumeAudio(buf coremedia.CMSampleBuffer) error {
	d.audioSamplesRcv++
	d.audioBytesRcv += uint64(len(buf.SampleData))
	return nil
}

func (d *DiagnosticsConsumer) consumeVideo(buf coremedia.CMSampleBuffer) error {
	d.videoSamplesRcv++
	d.videoBytesRcv += uint64(len(buf.SampleData))
	return nil
}
