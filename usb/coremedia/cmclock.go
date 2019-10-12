package coremedia

import (
	"time"
)

type CMClock struct {
	ID        uint64
	TimeScale uint32
	startTime time.Time
}

func NewCMClockWithHostTime(ID uint64) CMClock {
	return CMClock{
		ID: ID,
		//NanoSecond Scale
		TimeScale: 1000000000,
		startTime: time.Now(),
	}
}

func (c CMClock) GetTime() CMTime {
	return CMTime{
		CMTimeValue: uint64(time.Since(c.startTime).Nanoseconds()),
		CMTimeScale: c.TimeScale,
		CMTimeFlags: KCMTimeFlagsHasBeenRounded,
		CMTimeEpoch: 0,
	}

}
