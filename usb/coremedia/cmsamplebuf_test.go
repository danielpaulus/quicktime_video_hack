package coremedia_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/coremedia"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestCMSampleBuffer(t *testing.T) {
	dat, err := ioutil.ReadFile("../packet/fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}
	sbufPacket, err := coremedia.NewCMSampleBufferFromBytes(dat[20:])

	if assert.NoError(t, err) {
		assert.Equal(t, coremedia.KCMTimeFlags_HasBeenRounded, sbufPacket.OutputPresentationTimestamp.CMTimeFlags)
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
}
