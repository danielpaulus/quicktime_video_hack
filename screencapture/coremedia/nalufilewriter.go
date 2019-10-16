package coremedia

import (
	"encoding/binary"
	"io"
)

var startCode = []byte{00, 00, 00, 01}

//NaluFileWriter writes nalus into a file using 0x00000001 as a separator (h264 ANNEX B)
//The file is playable with vlc
type NaluFileWriter struct {
	outFileWriter io.Writer
	outFilePath   string
}

//NewNaluFileWriter binary writes nalus in annex b format to the given writer
func NewNaluFileWriter(outFileWriter io.Writer) NaluFileWriter {
	return NaluFileWriter{outFileWriter: outFileWriter}
}

//Consume writes PPS and SPS as well as sample bufs into a annex b .h264 file
func (nfw NaluFileWriter) Consume(buf CMSampleBuffer) error {
	if buf.HasFormatDescription {
		err := nfw.writeNalu(buf.FormatDescription.PPS)
		if err != nil {
			return err
		}
		err = nfw.writeNalu(buf.FormatDescription.SPS)
		if err != nil {
			return err
		}
	}
	return nfw.writeNalus(buf.SampleData)
}

func (nfw NaluFileWriter) writeNalus(bytes []byte) error {
	slice := bytes
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)
		err := nfw.writeNalu(slice[4 : length+4])
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (nfw NaluFileWriter) writeNalu(naluBytes []byte) error {
	_, err := nfw.outFileWriter.Write(startCode)
	if err != nil {
		return err
	}
	_, err = nfw.outFileWriter.Write(naluBytes)
	if err != nil {
		return err
	}
	return nil
}
