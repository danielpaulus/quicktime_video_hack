package coremedia_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/coremedia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeconds(t *testing.T) {
	time := createCmTime()
	assert.Equal(t, uint64(2), time.Seconds())
}

func TestErrors(t *testing.T) {
	time := createCmTime()
	buffer := make([]byte, 0)
	err := time.Serialize(buffer)
	assert.Error(t, err)
}

func TestCodec(t *testing.T) {
	time := createCmTime()
	buffer := make([]byte, 24)
	err := time.Serialize(buffer)
	if assert.NoError(t, err) {
		decodedTime, err := coremedia.NewCMTimeFromBytes(buffer)
		if assert.NoError(t, err) {
			assert.Equal(t, time.CMTimeEpoch, decodedTime.CMTimeEpoch)
			assert.Equal(t, time.CMTimeFlags, decodedTime.CMTimeFlags)
			assert.Equal(t, time.CMTimeScale, decodedTime.CMTimeScale)
			assert.Equal(t, time.CMTimeValue, decodedTime.CMTimeValue)
		}
	}
	_, err = coremedia.NewCMTimeFromBytes(buffer[:8])
	assert.Error(t, err)
}

func TestString(t *testing.T) {
	time := createCmTime()
	expected := "CMTime{1000/500, flags:KCMTimeFlagsHasBeenRounded, epoch:6}"
	s := time.String()
	assert.Equal(t, expected, s)
}

func createCmTime() coremedia.CMTime {
	return coremedia.CMTime{
		CMTimeValue: 1000,
		CMTimeScale: 500,
		CMTimeFlags: coremedia.KCMTimeFlagsHasBeenRounded,
		CMTimeEpoch: 6,
	}
}
