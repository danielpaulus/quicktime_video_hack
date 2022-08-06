package udpsink

import (
	"bytes"
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type udpsink struct {
	videoConn io.Writer
	audioConn io.Writer
	killFunc  func()
}

func New(video string, audio string) *udpsink {
	videoConn, err := net.Dial("tcp", video)
	if err != nil {
		log.Error(err)
		return nil
	}

	audioConn, err := net.Dial("tcp", audio)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &udpsink{videoConn, audioConn, func() {
		videoConn.Close()
		audioConn.Close()
	}}
}

func (u udpsink) sendWavHeader() {
	wavData, _ := coremedia.GetWavHeaderBytes(100)
	u.audioConn.Write(wavData)
}

var first bool = true

func (u udpsink) Consume(buf coremedia.CMSampleBuffer) error {
	if first {
		first = false
		u.sendWavHeader()

	}
	if buf.MediaType == coremedia.MediaTypeSound {

		buffer := make([]byte, 5000)
		reader := bytes.NewReader(buf.SampleData)
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			u.audioConn.Write(buffer[:n])
		}
		return nil
	}
	if buf.HasFormatDescription {
		log.Infof("%+v", buf.FormatDescription)
		//see above comment
		buf.OutputPresentationTimestamp.CMTimeValue = 0
		err := u.writeNalu(prependMarker(buf.FormatDescription.PPS, uint32(len(buf.FormatDescription.PPS))), buf)
		if err != nil {
			return err
		}
		err = u.writeNalu(prependMarker(buf.FormatDescription.SPS, uint32(len(buf.FormatDescription.SPS))), buf)
		if err != nil {
			return err
		}
	}
	u.writeNalus(buf)
	return nil
}
func (u udpsink) writeNalus(bytes coremedia.CMSampleBuffer) error {
	slice := bytes.SampleData
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)

		nalu := slice[4 : length+4]

		err := u.writeNalu(prependMarker(nalu, length), bytes)
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (u udpsink) writeNalu(naluBytes []byte, buf coremedia.CMSampleBuffer) error {
	buffer := make([]byte, 5000)
	reader := bytes.NewReader(naluBytes)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		u.videoConn.Write(buffer[:n])
	}

	return nil
}

var naluAnnexBMarkerBytes = []byte{0, 0, 0, 1}

func prependMarker(nalu []byte, length uint32) []byte {
	naluWithAnnexBMarker := make([]byte, length+4)
	copy(naluWithAnnexBMarker, naluAnnexBMarkerBytes)
	copy(naluWithAnnexBMarker[4:], nalu)
	return naluWithAnnexBMarker
}
func (u udpsink) Stop() {}
