package coremedia

import (
	"time"
)

// CMClock represents a monotonic Clock that will start counting when created
type CMClock struct {
	ID        uint64
	TimeScale uint32
	factor    float64
	startTime time.Time
}

//NanoSecondScale is the default system clock scale where 1/NanoSecondScale == 1 Nanosecond.
const NanoSecondScale = 1000000000

//NewCMClockWithHostTime creates a new Clock with the given ID with a nanosecond scale.
//Calls to GetTime will measure the time difference since the clock was created.
func NewCMClockWithHostTime(ID uint64) CMClock {
	return CMClock{
		ID:        ID,
		TimeScale: NanoSecondScale,
		factor:    1,
		startTime: time.Now(),
	}
}

//NewCMClockWithHostTimeAndScale creates a new CMClock with given ID and a custom timeScale
func NewCMClockWithHostTimeAndScale(ID uint64, timeScale uint32) CMClock {
	return CMClock{
		ID:        ID,
		TimeScale: timeScale,
		factor:    float64(timeScale) / float64(NanoSecondScale),
		startTime: time.Now(),
	}
}

//GetTime returns a CMTime that gives the time passed since the clock started.
//This is monotonic and does NOT use wallclock time.
func (c CMClock) GetTime() CMTime {
	return CMTime{
		CMTimeValue: c.calcValue(time.Since(c.startTime).Nanoseconds()),
		CMTimeScale: c.TimeScale,
		CMTimeFlags: KCMTimeFlagsHasBeenRounded,
		CMTimeEpoch: 0,
	}
}

func (c CMClock) calcValue(val int64) uint64 {
	if NanoSecondScale == c.TimeScale {
		return uint64(val)
	}
	return uint64(c.factor * float64(val))
}

//CalculateSkew calculates the deviation between the frequencies of two given clocks by using time diffs and returns a skew value float64
//scaled to match the second clock.
func CalculateSkew(startTimeClock1 CMTime, endTimeClock1 CMTime, startTimeClock2 CMTime, endTimeClock2 CMTime) float64 {
	timeDiffClock1 := endTimeClock1.CMTimeValue - startTimeClock1.CMTimeValue
	timeDiffClock2 := endTimeClock2.CMTimeValue - startTimeClock2.CMTimeValue

	diffTime := CMTime{CMTimeValue: timeDiffClock1, CMTimeScale: startTimeClock1.CMTimeScale}
	scaledDiff := diffTime.GetTimeForScale(startTimeClock2)
	//println("scaleddiff:" + scaledDiff)
	return float64(startTimeClock2.CMTimeScale) * scaledDiff / float64(timeDiffClock2)
}
