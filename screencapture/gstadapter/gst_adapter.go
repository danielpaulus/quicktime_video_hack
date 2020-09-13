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
	pipeline         *gst.Pipeline
	firstAudioSample bool
}

const audioAppSrcTargetElementName = "audio_target"
const videoAppSrcTargetElementName = "video_target"

const MP3 = "mp3"
const OGG = "ogg"

//New creates a new MAC OSX compatible gstreamer pipeline that will play device video and audio
//in a nice little window :-D
func New() *GstAdapter {
	log.Info("Starting Gstreamer..")
	pl := gst.NewPipeline("QT_Hack_Pipeline")

	videoAppSrc := setUpVideoPipeline(pl)
	audioAppSrc := setUpAudioPipelineBase(pl)
	setupLivePlayAudio(pl)

	pl.SetState(gst.STATE_PLAYING)
	runGlibMainLoop()

	log.Info("Gstreamer is running!")
	gsta := GstAdapter{videoAppSrc: videoAppSrc, audioAppSrc: audioAppSrc, firstAudioSample: true}

	return &gsta
}

func NewWithAudioPipeline(outfile string, audiotype string) (*GstAdapter, error) {
	log.Info("Starting Gstreamer..")
	pl := gst.NewPipeline("QT_Hack_Pipeline")

	audioAppSrc := setUpAudioPipelineBase(pl)
	switch audiotype {
	case MP3:
		setupMp3(pl, outfile)
	case OGG:
		setupVorbis(pl, outfile)
	default:
		log.Fatalf("Unrecognized Audio type:%s", audiotype)
	}

	pl.SetState(gst.STATE_PLAYING)
	runGlibMainLoop()

	log.Info("Gstreamer is running!")
	gsta := GstAdapter{audioAppSrc: audioAppSrc, firstAudioSample: true}

	return &gsta, nil
}

//NewWithCustomPipeline will parse the given pipelineString, connect the videoAppSrc to whatever element has the name "video_target" and the audioAppSrc to "audio_target"
//see also: https://gstreamer.freedesktop.org/documentation/application-development/appendix/programs.html?gi-language=c
func NewWithCustomPipeline(pipelineString string) (*GstAdapter, error) {
	log.Info("Starting Gstreamer..")
	log.WithFields(log.Fields{"custom_pipeline": pipelineString}).Debug("Starting Gstreamer with custom pipeline")
	pipeline, err := gst.ParseLaunch(pipelineString)
	if err != nil {
		return nil, fmt.Errorf("Invalid Pipeline, checkout --examples for help. Gstreamer parsing error was: %s", err)
	}

	audioAppSrcTargetElement := pipeline.AsBin().GetByName(audioAppSrcTargetElementName)
	if audioAppSrcTargetElement == nil {
		return nil, fmt.Errorf("The pipeline needs an element with a property 'name=%s' so I can link the audio source to it. run with --examples for details.", audioAppSrcTargetElementName)
	}

	videoAppSrcTargetElement := pipeline.AsBin().GetByName(videoAppSrcTargetElementName)
	if videoAppSrcTargetElement == nil {
		return nil, fmt.Errorf("The pipeline needs an element with a property 'name=%s' so I can link the video source to it. run with --examples for details.", videoAppSrcTargetElementName)
	}

	videoAppSrc := gst.NewAppSrc("my-video-src")
	videoAppSrc.SetProperty("is-live", true)

	audioAppSrc := gst.NewAppSrc("my-audio-src")
	audioAppSrc.SetProperty("is-live", true)

	pipeline.Add(videoAppSrc.AsElement())
	pipeline.Add(audioAppSrc.AsElement())

	audioAppSrc.Link(audioAppSrcTargetElement)
	videoAppSrc.Link(videoAppSrcTargetElement)

	pipeline.SetState(gst.STATE_PLAYING)
	//runGlibMainLoop()

	log.Info("Gstreamer is running!")
	gsta := GstAdapter{videoAppSrc: videoAppSrc, audioAppSrc: audioAppSrc, firstAudioSample: true, pipeline: pipeline}

	return &gsta, nil
}

//Stop sends an EOS (end of stream) event downstream the gstreamer pipeline.
//Some Elements need this to correctly finish. F.ex. writing mp4 video without
//sending EOS will result in a broken mp4 file
func (gsta GstAdapter) Stop() {
	log.Info("Stopping Gstreamer..")
	if gsta.audioAppSrc != nil {
		success := gsta.audioAppSrc.SendEvent(gst.Eos())
		if !success {
			log.Warn("Failed sending EOS signal for audio app source")
		}
	}
	if gsta.videoAppSrc != nil {
		success := gsta.videoAppSrc.SendEvent(gst.Eos())
		if !success {
			log.Warn("Failed sending EOS signal for video app source")
		}
	}

	if gsta.pipeline == nil {
		return
	}
	bus := gsta.pipeline.GetBus()

	//I hope those are 60 seconds
	msg := bus.TimedPopFiltered(1000000000*1000*60, gst.MESSAGE_EOS|gst.MESSAGE_ERROR)
	if msg == nil {
		log.Warn("No EOS received, video files might be broken")
		return
	}
	if msg.GetType() == gst.MESSAGE_ERROR {
		log.Warn("Error received, video files might be broken")
		return
	}
	log.Info("EOS received")
	gsta.pipeline.SetState(gst.STATE_NULL)
	log.Info("Gstreamer finished")
}

//runGlibMainLoop starts the glib Mainloop necessary for the video player to work on MAC OS X.
func runGlibMainLoop() {
	go func() {
		//See: https://golang.org/pkg/runtime/#LockOSThread
		runtime.LockOSThread()
		glib.NewMainLoop(nil).Run()
	}()
}
func setUpAudioPipelineBase(pl *gst.Pipeline) *gst.AppSrc {
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

	pl.Add(asrc.AsElement(), queue1, wavparse, audioconvert, queue2)
	asrc.Link(queue1)
	queue1.Link(wavparse)
	wavparse.Link(audioconvert)

	audioconvert.Link(queue2)

	return asrc
}
func setupVorbis(pl *gst.Pipeline, filepath string) {
	//vorbisenc ! oggmux ! filesink location=alsasrc.ogg
	vorbisEnc := gst.ElementFactoryMake("vorbisenc", "vorbisenc_01")
	checkElem(vorbisEnc, "vorbisenc_01")
	oggMux := gst.ElementFactoryMake("oggmux", "oggmux_01")
	checkElem(oggMux, "oggmux_01")

	filesink := gst.ElementFactoryMake("filesink", "filesink_01")
	filesink.SetProperty("location", filepath)
	checkElem(filesink, "filesink_01")

	pl.Add(vorbisEnc, oggMux, filesink)

	pl.GetByName("queue2").Link(vorbisEnc)
	vorbisEnc.Link(oggMux)
	oggMux.Link(filesink)
}
func setupMp3(pl *gst.Pipeline, filepath string) {
	// lamemp3enc ! filesink location=sine.mp3
	lameEnc := gst.ElementFactoryMake("lamemp3enc", "lamemp3enc_01")
	checkElem(lameEnc, "lamemp3enc_01")

	filesink := gst.ElementFactoryMake("filesink", "filesink_01")
	filesink.SetProperty("location", filepath)
	checkElem(filesink, "filesink_01")
	pl.Add(lameEnc, filesink)
	pl.GetByName("queue2").Link(lameEnc)
	lameEnc.Link(filesink)
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
