package rtpsupport

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

/*
gst-launch-1.0 -v udpsrc port=4001 caps='application/x-rtp, media=(string)audio, clock-rate=(int)48000, encoding-name=(string)L16, encoding-params=(string)1, channels=(int)2, payload=(int)96' ! rtpjitterbuffer latency=1000 ! rtpL16depay ! audioconvert ! autoaudiosink sync=true
*/
func TestNewRtpAudio(t *testing.T) {
	srv := NewRtpServer("localhost", 4000)
	srv.StartServerSocket()
	wavdata, err := ioutil.ReadFile("../../log/dump.wav")
	if err != nil {
		log.Fatal(err)
	}

	wavdata = wavdata[44:]
	for len(wavdata) > 1024 {
		samples := wavdata[:1024]
		wavdata = wavdata[1024:]
		beWav := make([]byte, 1024)
		for i := 0; i < 256; i++ {

			theint := binary.LittleEndian.Uint32(samples[i*4 : i*4+4])
			binary.BigEndian.PutUint32(beWav[i*4:i*4+4], theint)
		}

		sbuf := coremedia.CMSampleBuffer{
			SampleData: beWav,
			NumSamples: 256,
			MediaType:  coremedia.MediaTypeSound,
		}

		time.Sleep(time.Millisecond * 100)
		srv.Consume(sbuf)

	}
	log.Fatal("bla")
}

/*
func TestNewRtpServer(t *testing.T) {
	srv := NewRtpServer("localhost", 4000)
	var process *os.Process
	go func() {
		process = startGst()
		//process.Wait()
	}()
	srv.StartServerSocket()

	nalus, err := ioutil.ReadFile("/home/ganjalf/tmp/out.h264")
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
		time.Sleep(time.Millisecond * 10)
		timer.CMTimeValue += 16666666

	}
	time.Sleep(time.Second * 60)
	process.Kill()
}*/

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
