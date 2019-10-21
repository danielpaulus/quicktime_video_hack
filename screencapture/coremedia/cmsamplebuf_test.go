package coremedia_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

func TestCMSampleBuffer(t *testing.T) {
	dat, err := ioutil.ReadFile("../packet/fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}
	sbufPacket, err := coremedia.NewCMSampleBufferFromBytesVideo(dat[20:])

	if assert.NoError(t, err) {
		assert.Equal(t, true, sbufPacket.HasFormatDescription)
		assert.Equal(t, coremedia.KCMTimeFlagsHasBeenRounded, sbufPacket.OutputPresentationTimestamp.CMTimeFlags)
		assert.Equal(t, uint64(0x176a7), sbufPacket.OutputPresentationTimestamp.Seconds())
		assert.Equal(t, 1, len(sbufPacket.SampleTimingInfoArray))
		assert.Equal(t, uint64(0), sbufPacket.SampleTimingInfoArray[0].Duration.Seconds())
		assert.Equal(t, uint64(0x176a7), sbufPacket.SampleTimingInfoArray[0].PresentationTimeStamp.Seconds())
		assert.Equal(t, uint64(0), sbufPacket.SampleTimingInfoArray[0].DecodeTimeStamp.Seconds())
		assert.Equal(t, 90750, len(sbufPacket.SampleData))
		assert.Equal(t, 1, sbufPacket.NumSamples)
		assert.Equal(t, 1, len(sbufPacket.SampleSizes))
		assert.Equal(t, 90750, sbufPacket.SampleSizes[0])
		assert.Equal(t, 4, len(sbufPacket.Attachments.Entries))
		assert.Equal(t, 1, len(sbufPacket.Sary.Entries))
	}
	print(sbufPacket.String())
}

func TestCMSampleBufferNoFdsc(t *testing.T) {
	dat, err := ioutil.ReadFile("../packet/fixtures/asyn-feed-nofdsc")
	if err != nil {
		log.Fatal(err)
	}
	sbufPacket, err := coremedia.NewCMSampleBufferFromBytesVideo(dat[16:])

	if assert.NoError(t, err) {
		assert.Equal(t, false, sbufPacket.HasFormatDescription)
		assert.Equal(t, coremedia.KCMTimeFlagsHasBeenRounded, sbufPacket.OutputPresentationTimestamp.CMTimeFlags)
		assert.Equal(t, uint64(0x44b82fa09), sbufPacket.OutputPresentationTimestamp.Seconds())
		assert.Equal(t, 1, len(sbufPacket.SampleTimingInfoArray))
		assert.Equal(t, uint64(0), sbufPacket.SampleTimingInfoArray[0].Duration.Seconds())
		assert.Equal(t, uint64(0x44b82fa09), sbufPacket.SampleTimingInfoArray[0].PresentationTimeStamp.Seconds())
		assert.Equal(t, uint64(0), sbufPacket.SampleTimingInfoArray[0].DecodeTimeStamp.Seconds())
		assert.Equal(t, 56604, len(sbufPacket.SampleData))
		assert.Equal(t, 1, sbufPacket.NumSamples)
		assert.Equal(t, 1, len(sbufPacket.SampleSizes))
		assert.Equal(t, 56604, sbufPacket.SampleSizes[0])
		assert.Equal(t, 4, len(sbufPacket.Attachments.Entries))
		assert.Equal(t, 2, len(sbufPacket.Sary.Entries))
	}
	print(sbufPacket.String())
}

func TestCMSampleBufferAudio(t *testing.T) {
	dat, err := ioutil.ReadFile("../packet/fixtures/asyn-eat")
	if err != nil {
		log.Fatal(err)
	}
	sbufPacket, err := coremedia.NewCMSampleBufferFromBytesAudio(dat[16:])

	if assert.NoError(t, err) {
		assert.Equal(t, true, sbufPacket.HasFormatDescription)
		assert.Equal(t, 1024, sbufPacket.NumSamples)
		assert.Equal(t, 1, len(sbufPacket.SampleSizes))
		assert.Equal(t, 4, sbufPacket.SampleSizes[0])
		assert.Equal(t, sbufPacket.NumSamples*sbufPacket.SampleSizes[0], len(sbufPacket.SampleData))
		stringOutput := "{OutputPresentationTS:CMTime{2056/48000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, NumSamples:1024, SampleSize:4, fdsc:fdsc:{MediaType:Sound, AudioStreamBasicDescription: {SampleRate:48000.000000,FormatFlags:76,BytesPerPacket:4,FramesPerPacket:1,BytesPerFrame:4,ChannelsPerFrame:2,BitsPerChannel:16,Reserved:0}}}"
		assert.Equal(t, stringOutput, sbufPacket.String())
	}
	print(sbufPacket.String())
}

func TestCMSampleBufferAudioNoFdsc(t *testing.T) {
	dat, err := ioutil.ReadFile("../packet/fixtures/asyn-eat-nofdsc")
	if err != nil {
		log.Fatal(err)
	}
	sbufPacket, err := coremedia.NewCMSampleBufferFromBytesAudio(dat[16:])

	if assert.NoError(t, err) {
		assert.Equal(t, false, sbufPacket.HasFormatDescription)
		assert.Equal(t, 1024, sbufPacket.NumSamples)
		assert.Equal(t, 1, len(sbufPacket.SampleSizes))
		assert.Equal(t, 4, sbufPacket.SampleSizes[0])
		assert.Equal(t, sbufPacket.NumSamples*sbufPacket.SampleSizes[0], len(sbufPacket.SampleData))
		stringOutput := "{OutputPresentationTS:CMTime{3076/48000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, NumSamples:1024, SampleSize:4, fdsc:none}"
		assert.Equal(t, stringOutput, sbufPacket.String())
	}
}
