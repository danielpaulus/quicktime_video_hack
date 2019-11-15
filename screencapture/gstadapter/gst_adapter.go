package gstadapter

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/danielpaulus/gst"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/lijo-jose/glib"
	log "github.com/sirupsen/logrus"
)

//GstAdapter contains the AppSrc for accessing Gstreamer.
//TODO: add support for audio
//TODO: add support for shutting down gstreamer
//TODO: maybe add some queues to gst for performance enhancements
type GstAdapter struct {
	videoAppSrc      *gst.AppSrc
	audioAppSrc      *gst.AppSrc
	firstAudioSample bool
}

func Testvid() {
	src := gst.ElementFactoryMake("videotestsrc", "VideoSrc")
	checkElem(src, "videotestsrc")
	vsink := "autovideosink"

	sink := gst.ElementFactoryMake(vsink, "VideoSink")
	checkElem(sink, vsink)

	pl := gst.NewPipeline("MyPipeline")

	pl.Add(src, sink)

	src.Link(sink)
	pl.SetState(gst.STATE_PLAYING)

	glib.NewMainLoop(nil).Run()
}

func New() *GstAdapter {
	log.Info("Starting Gstreamer..")
	//Testvid()
	pl := gst.NewPipeline("QT_Hack_Pipeline")

	videoAppSrc := setUpVideoPipeline(pl)
	audioAppSrc := setUpAudioPipeline(pl)

	pl.SetState(gst.STATE_PLAYING)
	go func() { glib.NewMainLoop(nil).Run() }()
	//glib.NewMainLoop(nil).Run()
	log.Info("Gstreamer is running!")
	gsta := GstAdapter{videoAppSrc: videoAppSrc, audioAppSrc: audioAppSrc, firstAudioSample: true}
	//gsta := GstAdapter{videoAppSrc: videoAppSrc, firstAudioSample: true}
	return &gsta
}

func setUpAudioPipeline(pl *gst.Pipeline) *gst.AppSrc {
	asrc := gst.NewAppSrc("my-audio-src")
	asrc.SetProperty("is-live", true)

	filesink := gst.ElementFactoryMake("filesink", "filesink")
	checkElem(filesink, "filesink")
	filesink.SetProperty("location", "/home/ganjalf/tmp/audiodump.ogg")

	queue1 := gst.ElementFactoryMake("queue", "queue1")
	checkElem(queue1, "queue1")
	/*
		rawaudioparse := gst.ElementFactoryMake("rawaudioparse", "rawaudioparse_01")
		checkElem(rawaudioparse, "rawaudioparse_01")
		rawaudioparse.SetProperty("use-sink-caps", false)
		rawaudioparse.SetProperty("format", "pcm")
		rawaudioparse.SetProperty("pcm-format", "s16le")
		rawaudioparse.SetProperty("sample-rate", 48000)
		rawaudioparse.SetProperty("num-channels", 2)
	*/
	wavparse := gst.ElementFactoryMake("wavparse", "wavparse_01")
	checkElem(wavparse, "wavparse")
	wavparse.SetProperty("ignore-length", true)

	audioconvert := gst.ElementFactoryMake("audioconvert", "audioconvert_01")
	checkElem(audioconvert, "audioconvert_01")

	audioresample := gst.ElementFactoryMake("audioresample", "audioresample_01")
	checkElem(audioresample, "audioresample_01")

	autoaudiosink := gst.ElementFactoryMake("autoaudiosink", "autoaudiosink_01")
	checkElem(autoaudiosink, "autoaudiosink_01")
	autoaudiosink.SetProperty("sync", false)

	vorbisenc := gst.ElementFactoryMake("vorbisenc", "vorbisenc_01")
	checkElem(vorbisenc, "vorbisenc_01")

	oggmux := gst.ElementFactoryMake("oggmux", "oggmux_01")
	checkElem(oggmux, "oggmux_01")
	//vorbisenc ! oggmux ! filesink location=alsasrc.ogg

	//hack  oggdemux ! vorbisdec ! audioconvert

	oggdemux := gst.ElementFactoryMake("oggdemux", "oggdemux")
	checkElem(oggdemux, "oggdemux")

	vorbisdec := gst.ElementFactoryMake("vorbisdec", "vorbisdec")
	checkElem(vorbisdec, "vorbisdec")

	audioconvert2 := gst.ElementFactoryMake("audioconvert", "audioconvert_02")
	checkElem(audioconvert2, "audioconvert_02")

	//endhack

	pl.Add(asrc.AsElement(), queue1, wavparse, audioconvert, vorbisenc, oggmux, oggdemux, vorbisdec, audioconvert2, autoaudiosink)
	asrc.Link(queue1)
	queue1.Link(wavparse)
	wavparse.Link(audioconvert)

	audioconvert.Link(vorbisenc)

	vorbisenc.Link(vorbisdec)

	vorbisdec.Link(audioconvert2)
	audioconvert2.Link(autoaudiosink)
	//audioresample.Link(autoaudiosink)

	return asrc
}

func setUpVideoPipeline(pl *gst.Pipeline) *gst.AppSrc {
	asrc := gst.NewAppSrc("my-video-src")
	asrc.SetProperty("is-live", true)

	queue1 := gst.ElementFactoryMake("queue", "queue_11")
	checkElem(queue1, "queue_11")

	h264parse := gst.ElementFactoryMake("h264parse", "h264parse_01")
	checkElem(h264parse, "h264parse")

	avdec_h264 := gst.ElementFactoryMake("vaapih264dec", "avdec_h264_01")
	checkElem(avdec_h264, "avdec_h264_01")

	queue2 := gst.ElementFactoryMake("queue", "queue_12")
	checkElem(queue2, "queue_12")

	videoconvert := gst.ElementFactoryMake("videoconvert", "videoconvert_01")
	checkElem(videoconvert, "videoconvert_01")

	queue3 := gst.ElementFactoryMake("queue", "queue_13")
	checkElem(queue3, "queue_13")

	/*
		sink := gst.ElementFactoryMake("xvimagesink", "xvimagesink_01")
		checkElem(sink, "xvimagesink01")
	*/
	sink := gst.ElementFactoryMake("autovideosink", "autovideosink_01")
	//sink.SetProperty("sync", "false") does not do much
	checkElem(sink, "autovideosink_01")

	/*sink = gst.ElementFactoryMake("filesink", "filesink")
	sink.SetProperty("location", "/Users/danielpaulus/tmp/out-daniel.dump")
	checkElem(sink, "filesink")
	*/
	pl.Add(asrc.AsElement(), queue1, h264parse, avdec_h264, queue2, videoconvert, queue3, sink)

	asrc.Link(queue1)
	queue1.Link(h264parse)
	h264parse.Link(avdec_h264)
	avdec_h264.Link(queue2)
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

	if buf.OutputPresentationTimestamp.CMTimeValue > 17446044073700192000 {
		buf.OutputPresentationTimestamp.CMTimeValue = 0
	}
	if buf.HasFormatDescription {
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

func (nfw GstAdapter) writeNalus(bytes coremedia.CMSampleBuffer) error {
	slice := bytes.SampleData
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)

		nalu := slice[4 : length+4]

		err := nfw.writeNalu(prependMarker(nalu, length), bytes)
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (srv GstAdapter) writeNalu(naluBytes []byte, buf coremedia.CMSampleBuffer) error {
	naluLength := uint(len(naluBytes))
	gstBuf := gst.NewBufferAllocate(naluLength)
	//log.Infof("val:%d", buf.OutputPresentationTimestamp.CMTimeValue)
	gstBuf.SetPTS(buf.OutputPresentationTimestamp.CMTimeValue)
	gstBuf.SetDTS(0)
	//TODO: create CGO function that provides offsets so we can delete prependMarker again
	gstBuf.FillWithGoSlice(naluBytes)
	srv.videoAppSrc.PushBuffer(gstBuf)
	return nil
}

var naluAnnexBMarkerBytes = []byte{0, 0, 0, 1}

func prependMarker(nalu []byte, length uint32) []byte {
	naluWithAnnexBMarker := make([]byte, length+4)
	copy(naluWithAnnexBMarker, naluAnnexBMarkerBytes)
	copy(naluWithAnnexBMarker[4:], nalu)
	return naluWithAnnexBMarker
}
