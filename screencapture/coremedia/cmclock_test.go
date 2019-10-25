package coremedia_test

import (
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
)

func TestCMClock_GetTime(t *testing.T) {
	cmclock := coremedia.NewCMClockWithHostTime(uint64(5))
	assert.Equal(t, uint32(coremedia.NanoSecondScale), cmclock.TimeScale)
	time1 := cmclock.GetTime()
	time2 := cmclock.GetTime()
	assert.Equal(t, cmclock.TimeScale, time1.CMTimeScale)
	//The clock is monotonic with nanosecond precision, so this should always be true
	assert.True(t, true, time2.CMTimeValue > time1.CMTimeValue)

	cmclock = coremedia.NewCMClockWithHostTimeAndScale(0, 1)
	assert.Equal(t, uint64(0), cmclock.GetTime().CMTimeValue)
}

func TestCalculateSkew(t *testing.T) {
	testCases := map[string]struct {
		startTimeClock1 coremedia.CMTime
		endTimeClock1   coremedia.CMTime
		startTimeClock2 coremedia.CMTime
		endTimeClock2   coremedia.CMTime
		expectedValue   float64
	}{
		"check simple case, no skew": {
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 48000},
			float64(48000.0)},
		"check simple case, positive skew": {
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 2, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 1, CMTimeScale: 48000},
			float64(96000.0)},
		"check simple case, negative skew": {
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 2000, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 2001, CMTimeScale: 48000},
			float64(47976.011994003)},
		"check different scales, negative skew": {
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 20833 * 5, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 5, CMTimeScale: 48000},
			float64(47999.232)},
		"check different scales, positive skew": {
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 20833 * 5001, CMTimeScale: coremedia.NanoSecondScale},
			coremedia.CMTime{CMTimeValue: 0, CMTimeScale: 48000},
			coremedia.CMTime{CMTimeValue: 5000, CMTimeScale: 48000},
			float64(48008.8318464)},
	}

	for s, tc := range testCases {
		calculatedSkew := coremedia.CalculateSkew(tc.startTimeClock1, tc.endTimeClock1, tc.startTimeClock2, tc.endTimeClock2)
		assert.Equal(t, tc.expectedValue, calculatedSkew, s)
	}
}
