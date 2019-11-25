// +build darwin

package gstadapter

import (
	"encoding/binary"
	"fmt"
	"os"
	"runtime"

	"github.com/danielpaulus/gst"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/lijo-jose/glib"
	log "github.com/sirupsen/logrus"
)

//GstAdapter contains the AppSrc for accessing Gstreamer.
type GstAdapter struct {
	videoAppSrc      *gst.AppSrc
	audioAppSrc      *gst.AppSrc
	firstAudioSample bool
}

//New creates a new MAC OSX compatible gstreamer pipeline that will play device video and audio
//in a nice little window :-D
func New() *GstAdapter {
	log.Info("Starting Gstreamer..")
	pl := gst.NewPipeline("QT_Hack_Pipeline")

	videoAppSrc := setUpVideoPipeline(pl)
	audioAppSrc := setUpAudioPipeline(pl)

	pl.SetState(gst.STATE_PLAYING)
	runGlibMainLoop()

	log.Info("Gstreamer is running!")
	gsta := GstAdapter{videoAppSrc: videoAppSrc, audioAppSrc: audioAppSrc, firstAudioSample: true}

	return &gsta
}

//runGlibMainLoop starts the glib Mainloop necessary for the video player to work on MAC OS X.
func runGlibMainLoop() {
	go func() {
		//See: https://golang.org/pkg/runtime/#LockOSThread
		runtime.LockOSThread()
		glib.NewMainLoop(nil).Run()
	}()
}

func setUpAudioPipeline(pl *gst.Pipeline) *gst.AppSrc {
	asrc := gst.NewAppSrc("my-audio-src")
	asrc.SetProperty("is-live", true)

	queue1 := gst.ElementFactoryMake("queue", "queue1")
	checkElem(queue1, "queue1")

	queue2 := gst.ElementFactoryMake("queue", "queue2")
	checkElem(queue1, "queue2")

	wavparse := gst.ElementFactoryMake("wavparse", "wavparse_01")
	checkElem(wavparse, "wavparse")
	wavparse.SetProperty("ignore-length", true)

	audioconvert := gst.ElementFactoryMake("audioconvert", "audioconvert_01")
	checkElem(audioconvert, "audioconvert_01")

	autoaudiosink := gst.ElementFactoryMake("autoaudiosink", "autoaudiosink_01")
	checkElem(autoaudiosink, "autoaudiosink_01")
	autoaudiosink.SetProperty("sync", false)

	pl.Add(asrc.AsElement(), queue1, wavparse, audioconvert, queue2, autoaudiosink)
	asrc.Link(queue1)
	queue1.Link(wavparse)
	wavparse.Link(audioconvert)

	audioconvert.Link(queue2)
	queue2.Link(autoaudiosink)

	return asrc
}

func setUpVideoPipeline(pl *gst.Pipeline) *gst.AppSrc {
	asrc := gst.NewAppSrc("my-video-src")
	asrc.SetProperty("is-live", true)

	queue1 := gst.ElementFactoryMake("queue", "queue_11")
	checkElem(queue1, "queue_11")

	h264parse := gst.ElementFactoryMake("h264parse", "h264parse_01")
	checkElem(h264parse, "h264parse")

	avdecH264 := gst.ElementFactoryMake("vtdec", "avdec_h264_01")
	checkElem(avdecH264, "avdec_h264_01")

	queue2 := gst.ElementFactoryMake("queue", "queue_12")
	checkElem(queue2, "queue_12")

	videoconvert := gst.ElementFactoryMake("videoconvert", "videoconvert_01")
	checkElem(videoconvert, "videoconvert_01")

	queue3 := gst.ElementFactoryMake("queue", "queue_13")
	checkElem(queue3, "queue_13")

	sink := gst.ElementFactoryMake("autovideosink", "autovideosink_01")
	// setting this to true, creates extremely choppy video
	// I probably messed up something regarding the time stamps
	sink.SetProperty("sync", false)
	checkElem(sink, "autovideosink_01")

	pl.Add(asrc.AsElement(), queue1, h264parse, avdecH264, queue2, videoconvert, queue3, sink)

	asrc.Link(queue1)
	queue1.Link(h264parse)
	h264parse.Link(avdecH264)
	avdecH264.Link(queue2)
	queue2.Link(videoconvert)
	videoconvert.Link(queue3)
	queue3.Link(sink)
	return asrc
}

func checkElem(e *gst.Element, name string) {
	if e == nil {
		fmt.Fprintln(os.Stderr, "can't make element: ", name)
		os.Exit(1)
	}
}

//Consume will transfer AV data into a Gstreamer AppSrc
func (gsta *GstAdapter) Consume(buf coremedia.CMSampleBuffer) error {
	if buf.MediaType == coremedia.MediaTypeSound {
		if gsta.firstAudioSample {
			gsta.firstAudioSample = false
			gsta.sendWavHeader()
		}
		return gsta.sendAudioSample(buf)
	}

	//FIXME: ugly hack I added to prevent gstreamer from receiving decreasing timestamps
	//I might have messed something up while sending times to the device as my first
	//buffer will have this weird, large timestamp. So I hack it to be equal to zero here
	if buf.OutputPresentationTimestamp.CMTimeValue > 17446044073700192000 {
		buf.OutputPresentationTimestamp.CMTimeValue = 0
	}
	if buf.HasFormatDescription {
		//see above comment
		buf.OutputPresentationTimestamp.CMTimeValue = 0
		err := gsta.writeNalu(prependMarker(buf.FormatDescription.PPS, uint32(len(buf.FormatDescription.PPS))), buf)
		if err != nil {
			return err
		}
		err = gsta.writeNalu(prependMarker(buf.FormatDescription.SPS, uint32(len(buf.FormatDescription.SPS))), buf)
		if err != nil {
			return err
		}
	}
	gsta.writeNalus(buf)

	return nil
}

func (gsta GstAdapter) sendWavHeader() {
	wavData, _ := coremedia.GetWavHeaderBytes(100)
	sampleLength := uint(len(wavData))
	gstBuf := gst.NewBufferAllocate(sampleLength)
	gstBuf.SetPTS(0)
	gstBuf.SetDTS(0)
	//TODO: create CGO function that provides offsets so we can delete prependMarker again
	gstBuf.FillWithGoSlice(wavData)
	gsta.audioAppSrc.PushBuffer(gstBuf)
}

func (gsta GstAdapter) sendAudioSample(buf coremedia.CMSampleBuffer) error {
	sampleLength := uint(len(buf.SampleData))
	gstBuf := gst.NewBufferAllocate(sampleLength)
	gstBuf.SetPTS(buf.OutputPresentationTimestamp.CMTimeValue)
	gstBuf.SetDTS(0)
	//TODO: create CGO function that provides offsets so we can delete prependMarker again
	gstBuf.FillWithGoSlice(buf.SampleData)
	gsta.audioAppSrc.PushBuffer(gstBuf)

	return nil
}

func (gsta GstAdapter) writeNalus(bytes coremedia.CMSampleBuffer) error {
	slice := bytes.SampleData
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)

		nalu := slice[4 : length+4]

		err := gsta.writeNalu(prependMarker(nalu, length), bytes)
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (gsta GstAdapter) writeNalu(naluBytes []byte, buf coremedia.CMSampleBuffer) error {
	naluLength := uint(len(naluBytes))
	gstBuf := gst.NewBufferAllocate(naluLength)

	gstBuf.SetPTS(buf.OutputPresentationTimestamp.CMTimeValue)
	gstBuf.SetDTS(0)
	//TODO: create CGO function that provides offsets so we can delete prependMarker again
	gstBuf.FillWithGoSlice(naluBytes)
	gsta.videoAppSrc.PushBuffer(gstBuf)
	return nil
}

var naluAnnexBMarkerBytes = []byte{0, 0, 0, 1}

func prependMarker(nalu []byte, length uint32) []byte {
	naluWithAnnexBMarker := make([]byte, length+4)
	copy(naluWithAnnexBMarker, naluAnnexBMarkerBytes)
	copy(naluWithAnnexBMarker[4:], nalu)
	return naluWithAnnexBMarker
}
