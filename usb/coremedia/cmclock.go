package coremedia

import (
	"time"
)

type CMClock struct {
	ID        uint64
	TimeScale uint32
	startTime time.Time
}

//NewCMClockWithHostTime creates a new Clock with the given ID with a nanosecond scale.
//Calls to GetTime will measure the time difference since the clock was created.
func NewCMClockWithHostTime(ID uint64) CMClock {
	return CMClock{
		ID: ID,
		//NanoSecond Scale
		TimeScale: 1000000000,
		startTime: time.Now(),
	}
}

//GetTime returns a CMTime that gives the time passed since the clock started.
//This is monotonic and does NOT use wallclock time.
func (c CMClock) GetTime() CMTime {
	return CMTime{
		CMTimeValue: uint64(time.Since(c.startTime).Nanoseconds()),
		CMTimeScale: c.TimeScale,
		CMTimeFlags: KCMTimeFlagsHasBeenRounded,
		CMTimeEpoch: 0,
	}

}
