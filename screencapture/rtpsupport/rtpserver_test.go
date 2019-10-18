package rtpsupport

import (
	"bytes"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNewRtpServer(t *testing.T) {
	srv := NewRtpServer()
	var process *os.Process
	go func() {
		process = startGst()
		//process.Wait()
	}()
	srv.StartServerSocket()

	nalus, err := ioutil.ReadFile("/home/ganjalf/out.h264")
	if err != nil {
		log.Fatal(err)
	}
	singleNalus := bytes.Split(nalus, []byte{0, 0, 0, 1})
	timer := coremedia.CMTime{
		CMTimeValue: 0,
		CMTimeScale: 1000000000,
		CMTimeFlags: 0,
		CMTimeEpoch: 0,
	}
	for _, nalu := range singleNalus {
		sbuf := coremedia.CMSampleBuffer{
			SampleData:                  nalu,
			OutputPresentationTimestamp: timer,
		}

		srv.Consume(sbuf)
		time.Sleep(time.Millisecond*10)
		timer.CMTimeValue += 16666666

	}
	time.Sleep(time.Second * 60)
	process.Kill()
}
func startGst() *os.Process {
	cmd := exec.Command("gst-launch-1.0", "-v", "udpsrc", "port=5000", "caps=\"application/x-rtp, media=(string)video, clock-rate=(int)90000, encoding-name=(string)H264, payload=(int)96\"", "!", "rtph264depay", "!",
		"decodebin", "!", "videoconvert", "!", "autovideosink")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	/*err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}*/
	return cmd.Process
}
