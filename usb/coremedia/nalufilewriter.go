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
func (nfw NaluFileWriter) Consume(buf CMSampleBuffer) {
	if buf.HasFormatDescription {
		nfw.writeNalu(buf.FormatDescription.PPS)
		nfw.writeNalu(buf.FormatDescription.SPS)
	}
	nfw.writeNalus(buf.SampleData)
}

func (nfw NaluFileWriter) writeNalus(bytes []byte) {
	slice := bytes
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)
		nfw.writeNalu(slice[4 : length+4])
		slice = slice[length+4:]
	}
}

func (nfw NaluFileWriter) writeNalu(naluBytes []byte) {
	nfw.outFileWriter.Write(startCode)
	nfw.outFileWriter.Write(naluBytes)
}
