package coremedia_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/dict"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakePPS = []byte{1, 2, 3, 4, 5}
var fakeSPS = []byte{6, 7, 8, 9}
var startCode = []byte{00, 00, 00, 01}

func TestFileWriter(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 100))
	buf.Reset()
	nfw := coremedia.NewNaluFileWriter(buf)
	err := nfw.Consume(cmSampleBufWithAFewBytes())
	assert.NoError(t, err)
	assert.Equal(t, 6, buf.Len())
	assert.Equal(t, []byte{00, 00, 00, 01, 00, 00}, buf.Bytes())
	buf.Reset()
	err = nfw.Consume(cmSampleBufWithFdscAndAFewBytes())
	assert.NoError(t, err)
	expectedBytes := append(startCode, fakePPS...)
	expectedBytes = append(expectedBytes, startCode...)
	expectedBytes = append(expectedBytes, fakeSPS...)
	expectedBytes = append(expectedBytes, []byte{00, 00, 00, 01, 00, 00}...)
	assert.Equal(t, expectedBytes, buf.Bytes())

	nfw = coremedia.NewNaluFileWriter(failingWriter{})
	err = nfw.Consume(cmSampleBufWithFdscAndAFewBytes())
	assert.Error(t, err)
	err = nfw.Consume(cmSampleBufWithAFewBytes())
	assert.Error(t, err)
}

type failingWriter struct{}

func (f failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("failed")
}

func cmSampleBufWithFdscAndAFewBytes() coremedia.CMSampleBuffer {
	fakeNalu := make([]byte, 6)
	binary.BigEndian.PutUint32(fakeNalu, 2)
	return cmSampleBufWithFdscAndSampleData(fakeNalu)
}

func cmSampleBufWithFdscAndSampleData(sampleData []byte) coremedia.CMSampleBuffer {
	return coremedia.CMSampleBuffer{
		OutputPresentationTimestamp: coremedia.CMTime{},
		FormatDescription:           dict.FormatDescriptor{PPS: fakePPS, SPS: fakeSPS},
		HasFormatDescription:        true,
		NumSamples:                  0,
		SampleTimingInfoArray:       nil,
		SampleData:                  sampleData,
		SampleSizes:                 nil,
		Attachments:                 dict.IndexKeyDict{},
		Sary:                        dict.IndexKeyDict{},
	}
}

func cmSampleBufWithAFewBytes() coremedia.CMSampleBuffer {
	fakeNalu := make([]byte, 6)
	binary.BigEndian.PutUint32(fakeNalu, 2)
	return cmSampleBufWithSampleData(fakeNalu)
}
func cmSampleBufWithSampleData(sampleData []byte) coremedia.CMSampleBuffer {
	return coremedia.CMSampleBuffer{
		OutputPresentationTimestamp: coremedia.CMTime{},
		FormatDescription:           dict.FormatDescriptor{},
		HasFormatDescription:        false,
		NumSamples:                  0,
		SampleTimingInfoArray:       nil,
		SampleData:                  sampleData,
		SampleSizes:                 nil,
		Attachments:                 dict.IndexKeyDict{},
		Sary:                        dict.IndexKeyDict{},
	}
}
