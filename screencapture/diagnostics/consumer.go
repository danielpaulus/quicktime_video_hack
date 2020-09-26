package diagnostics

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	log "github.com/sirupsen/logrus"
)

//CSVHeader contains the header for the metrics file
const CSVHeader = "audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv, heapobjects, alloc\n"

//DiagnosticsConsumer periodically logs samples received, bytes received and memory stats to a csv file.
type DiagnosticsConsumer struct {
	outFileWriter   io.Writer
	audioSamplesRcv uint64
	videoSamplesRcv uint64
	audioBytesRcv   uint64
	videoBytesRcv   uint64
	mux             sync.Mutex
	interval        time.Duration
	stop            chan struct{}
	stopDone        chan struct{}
}

//NewDiagnosticsConsumer creates a new DiagnosticsConsumer
func NewDiagnosticsConsumer(outfile io.Writer, interval time.Duration) *DiagnosticsConsumer {
	d := &DiagnosticsConsumer{outFileWriter: outfile, interval: interval, stop: make(chan struct{}), stopDone: make(chan struct{})}
	go fileWriter(d)
	return d
}

func fileWriter(d *DiagnosticsConsumer) {
	d.outFileWriter.Write([]byte(CSVHeader))

	for {

		select {
		case <-d.stop:
			log.Info("Stopped")
			close(d.stopDone)
			return
		case <-time.After(d.interval):
			audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv := readAndReset(d)
			heapobjects, alloc := getMemStats()
			csvLine := fmt.Sprintf("%d,%d,%d,%d,%d,%d\n", audioSamplesRcv, audioBytesRcv, videoSamplesRcv, videoBytesRcv, heapobjects, alloc)
			_, err := d.outFileWriter.Write([]byte(csvLine))
			if err != nil {
				log.Fatalf("Failed writing to metricsfile:%+v", err)
			}
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

//Consume logs stats
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

//Stop writing to the csv file
func (d *DiagnosticsConsumer) Stop() {
	close(d.stop)
	<-d.stopDone
}
