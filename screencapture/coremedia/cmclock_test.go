package coremedia_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCMClock_GetTime(t *testing.T) {
	cmclock := coremedia.NewCMClockWithHostTime(uint64(5))
	assert.Equal(t, uint32(1000000000), cmclock.TimeScale)
	time1 := cmclock.GetTime()
	time2 := cmclock.GetTime()
	assert.Equal(t, cmclock.TimeScale, time1.CMTimeScale)
	//The clock is monotonic with nanosecond precision, so this should always be true
	assert.True(t, true, time2.CMTimeValue > time1.CMTimeValue)
}
