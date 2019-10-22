package coremedia_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

var fakePPS = []byte{1, 2, 3, 4, 5}
var fakeSPS = []byte{6, 7, 8, 9}
var startCode = []byte{00, 00, 00, 01}

func TestFileWriter(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 100))
	buf.Reset()
	avfw := coremedia.NewAVFileWriter(buf, nil)
	err := avfw.Consume(cmSampleBufWithAFewBytes())
	assert.NoError(t, err)
	assert.Equal(t, 6, buf.Len())
	assert.Equal(t, []byte{00, 00, 00, 01, 00, 00}, buf.Bytes())
	buf.Reset()
	err = avfw.Consume(cmSampleBufWithFdscAndAFewBytes())
	assert.NoError(t, err)
	expectedBytes := append(startCode, fakePPS...)
	expectedBytes = append(expectedBytes, startCode...)
	expectedBytes = append(expectedBytes, fakeSPS...)
	expectedBytes = append(expectedBytes, []byte{00, 00, 00, 01, 00, 00}...)
	assert.Equal(t, expectedBytes, buf.Bytes())

	avfw = coremedia.NewAVFileWriter(failingWriter{}, nil)
	err = avfw.Consume(cmSampleBufWithFdscAndAFewBytes())
	assert.Error(t, err)
	err = avfw.Consume(cmSampleBufWithAFewBytes())
	assert.Error(t, err)
}

func TestFileWriterForAudio(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 100))
	buf.Reset()
	avfw := coremedia.NewAVFileWriter(nil, buf)
	sampleBuffer := cmSampleBufWithFdscAndAFewBytes()
	sampleBuffer.MediaType = coremedia.MediaTypeSound
	err := avfw.Consume(sampleBuffer)
	if assert.NoError(t, err) {
		assert.Equal(t, sampleBuffer.SampleData, buf.Bytes())
	}

	avfw = coremedia.NewAVFileWriter(nil, failingWriter{})
	err = avfw.Consume(sampleBuffer)
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
		FormatDescription:           coremedia.FormatDescriptor{PPS: fakePPS, SPS: fakeSPS},
		HasFormatDescription:        true,
		NumSamples:                  0,
		SampleTimingInfoArray:       nil,
		SampleData:                  sampleData,
		SampleSizes:                 nil,
		Attachments:                 coremedia.IndexKeyDict{},
		Sary:                        coremedia.IndexKeyDict{},
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
		FormatDescription:           coremedia.FormatDescriptor{},
		HasFormatDescription:        false,
		NumSamples:                  0,
		SampleTimingInfoArray:       nil,
		SampleData:                  sampleData,
		SampleSizes:                 nil,
		Attachments:                 coremedia.IndexKeyDict{},
		Sary:                        coremedia.IndexKeyDict{},
	}
}
