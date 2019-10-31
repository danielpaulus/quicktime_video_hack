package gstadapter

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/danielpaulus/gst"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	log "github.com/sirupsen/logrus"
)

//GstAdapter contains the AppSrc for accessing Gstreamer.
//TODO: add support for audio
//TODO: add support for shutting down gstreamer
//TODO: maybe add some queues to gst for performance enhancements
type GstAdapter struct {
	appSrc *gst.AppSrc
}

func New() *GstAdapter {
	log.Info("Starting Gstreamer..")
	asrc := gst.NewAppSrc("my-video-src")
	asrc.SetProperty("is-live", true)

	sink := gst.ElementFactoryMake("xvimagesink", "xvimagesink_01")
	checkElem(sink, "xvimagesink01")

	h264parse := gst.ElementFactoryMake("h264parse", "h264parse_01")
	checkElem(h264parse, "h264parse")

	avdec_h264 := gst.ElementFactoryMake("avdec_h264", "avdec_h264_01")
	checkElem(avdec_h264, "avdec_h264_01")

	videoconvert := gst.ElementFactoryMake("videoconvert", "videoconvert_01")
	checkElem(videoconvert, "videoconvert_01")

	pl := gst.NewPipeline("QT_Hack_Pipeline")
	pl.Add(asrc.AsElement(), h264parse, avdec_h264, videoconvert, sink)

	asrc.Link(h264parse)
	h264parse.Link(avdec_h264)
	avdec_h264.Link(videoconvert)
	videoconvert.Link(sink)

	pl.SetState(gst.STATE_PLAYING)

	log.Info("Gstreamer is running!")
	gsta := GstAdapter{appSrc: asrc}
	return &gsta
}

func checkElem(e *gst.Element, name string) {
	if e == nil {
		fmt.Fprintln(os.Stderr, "can't make element: ", name)
		os.Exit(1)
	}
}

//Consume will transfer AV data into a Gstreamer AppSrc
func (gsta GstAdapter) Consume(buf coremedia.CMSampleBuffer) error {
	if buf.MediaType == coremedia.MediaTypeSound {
		return gsta.sendAudioSample(buf)
	}

	if buf.HasFormatDescription {
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

func (gsta GstAdapter) sendAudioSample(buf coremedia.CMSampleBuffer) error {
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
	gstBuf.SetPTS(buf.OutputPresentationTimestamp.CMTimeValue)
	gstBuf.SetDTS(0)
	//TODO: create CGO function that provides offsets so we can delete prependMarker again
	gstBuf.FillWithGoSlice(naluBytes)
	srv.appSrc.PushBuffer(gstBuf)
	return nil
}

var naluAnnexBMarkerBytes = []byte{0, 0, 0, 1}

func prependMarker(nalu []byte, length uint32) []byte {
	naluWithAnnexBMarker := make([]byte, length+4)
	copy(naluWithAnnexBMarker, naluAnnexBMarkerBytes)
	copy(naluWithAnnexBMarker[4:], nalu)
	return naluWithAnnexBMarker
}
