package coremedia

import (
	"encoding/binary"
	"io"
)

var startCode = []byte{00, 00, 00, 01}

//AVFileWriter writes nalus into a file using 0x00000001 as a separator (h264 ANNEX B) and raw pcm audio into a wav file
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
type AVFileWriter struct {
	outFileWriter   io.Writer
	audioFileWriter io.Writer
	outFilePath     string
}

//NewAVFileWriter binary writes nalus in annex b format to the given writer and audio buffers into a wav file.
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
func NewAVFileWriter(outFileWriter io.Writer, audioFileWriter io.Writer) AVFileWriter {
	return AVFileWriter{outFileWriter: outFileWriter, audioFileWriter: audioFileWriter}
}

//Consume writes PPS and SPS as well as sample bufs into a annex b .h264 file and audio samples into a wav file
//Note that you will have to call WriteWavHeader() on the audiofile when you are done to write a wav header and get a valid file.
func (avfw AVFileWriter) Consume(buf CMSampleBuffer) error {
	if buf.MediaType == MediaTypeSound {
		return avfw.consumeAudio(buf)
	}
	return avfw.consumeVideo(buf)
}

func (avfw AVFileWriter) consumeVideo(buf CMSampleBuffer) error {
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
	return avfw.writeNalus(buf.SampleData)
}

func (avfw AVFileWriter) writeNalus(bytes []byte) error {
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

func (avfw AVFileWriter) writeNalu(naluBytes []byte) error {
	_, err := avfw.outFileWriter.Write(startCode)
	if err != nil {
		return err
	}
	_, err = avfw.outFileWriter.Write(naluBytes)
	if err != nil {
		return err
	}
	return nil
}

func (avfw AVFileWriter) consumeAudio(buffer CMSampleBuffer) error {
	_, err := avfw.audioFileWriter.Write(buffer.SampleData)
	if err != nil {
		return err
	}
	return nil
}
