package coremedia_test

import (
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

func TestScaleConversion(t *testing.T) {
	testCases := map[string]struct {
		originalTime         coremedia.CMTime
		destinationScaleTime coremedia.CMTime
		expectedTime         float64
	}{
		"check zero valued CMTime works": {coremedia.CMTime{CMTimeValue: 0, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 2 * coremedia.NanoSecondScale},
			0},
		"doubling scale": {coremedia.CMTime{CMTimeValue: 100, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 2 * coremedia.NanoSecondScale},
			float64(0xC8)},
		"smaller scale": {coremedia.CMTime{CMTimeValue: 100, CMTimeScale: 1},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 48000},
			float64(0x493e00)},
	}

	for s, tc := range testCases {
		actualTime := tc.originalTime.GetTimeForScale(tc.destinationScaleTime)
		assert.Equal(t, tc.expectedTime, actualTime, s)
	}
}

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
