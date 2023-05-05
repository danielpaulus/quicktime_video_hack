package gstadapter

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

var startCode = []byte{00, 00, 00, 01}

//AVFileWriter writes nalus into a file using 0x00000001 as a separator (h264 ANNEX B) and raw pcm audio into a wav file
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
type TCPServerWriter struct {
	h264FileWriter  io.Writer
	wavFileWriter   io.Writer
	outFilePath     string
	audioOnly       bool
	videoConnWaiter chan net.Conn
	audioConnWaiter chan net.Conn
	errReceiver     chan error
}

const (
	CONN_HOST       = "localhost"
	CONN_PORT       = "3333"
	AUDIO_CONN_PORT = "3334"
	CONN_TYPE       = "tcp"
)

//NewAVFileWriter binary writes nalus in annex b format to the given writer and audio buffers into a wav file.
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
func StartTcpWriter() (TCPServerWriter, error) {
	// Listen for incoming connections.
	videoConnWaiter := make(chan net.Conn)
	audioConnWaiter := make(chan net.Conn)
	errReceiver := make(chan error)
	go func() {
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			log.Error("Error listening:", err.Error())
			errReceiver <- err
			return
		}
		log.Info("waiting for video connection... on 3333")
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting: ", err.Error())
			errReceiver <- err
			return
		}
		log.Info("ok!")
		videoConnWaiter <- conn
	}()

	go func() {
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+AUDIO_CONN_PORT)
		if err != nil {
			log.Error("Error listening:", err.Error())
			errReceiver <- err
			return
		}
		log.Info("waiting for audio connection... on 3334")
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting: ", err.Error())
			errReceiver <- err
			return
		}
		log.Info("ok!")
		audioConnWaiter <- conn
	}()

	var wavWriter io.Writer
	var h264Writer io.Writer
	select {
	case ac := <-audioConnWaiter:
		wavWriter = ac
	case err := <-errReceiver:
		return TCPServerWriter{}, err
	}
	select {
	case vc := <-videoConnWaiter:
		h264Writer = vc
	case err := <-errReceiver:
		return TCPServerWriter{}, err
	}
	return TCPServerWriter{h264FileWriter: h264Writer, wavFileWriter: wavWriter, audioOnly: false}, nil
}

//Consume writes PPS and SPS as well as sample bufs into a annex b .h264 file and audio samples into a wav file
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
func (avfw TCPServerWriter) Consume(buf coremedia.CMSampleBuffer) error {
	if buf.MediaType == coremedia.MediaTypeSound {
		return avfw.consumeAudio(buf)
	}
	if avfw.audioOnly {
		return nil
	}
	return avfw.consumeVideo(buf)
}

//Nothing currently
func (avfw TCPServerWriter) Stop() {}

func (avfw TCPServerWriter) consumeVideo(buf coremedia.CMSampleBuffer) error {
	if buf.HasFormatDescription {
		err := avfw.writeNalu(buf.FormatDescription.PPS)
		if err != nil {
			return err
		}
		err = avfw.writeNalu(buf.FormatDescription.SPS)
		if err != nil {
			return err
		}
	}
	if !buf.HasSampleData() {
		return nil
	}
	return avfw.writeNalus(buf.SampleData)
}

func (avfw TCPServerWriter) writeNalus(bytes []byte) error {
	slice := bytes
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)
		err := avfw.writeNalu(slice[4 : length+4])
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (avfw TCPServerWriter) writeNalu(naluBytes []byte) error {
	_, err := avfw.h264FileWriter.Write(startCode)
	if err != nil {
		return err
	}
	_, err = avfw.h264FileWriter.Write(naluBytes)
	if err != nil {
		return err
	}
	return nil
}

func (avfw TCPServerWriter) consumeAudio(buffer coremedia.CMSampleBuffer) error {
	if !buffer.HasSampleData() {
		return nil
	}
	_, err := avfw.wavFileWriter.Write(buffer.SampleData)
	if err != nil {
		return err
	}
	return nil
}
