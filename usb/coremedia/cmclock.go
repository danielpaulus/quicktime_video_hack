package coremedia

import (
	"time"
)

type CMClock struct {
	ID        uint64
	TimeScale uint32
}

func (c CMClock) GetTime() CMTime {
	currentTime := time.Now().UnixNano()
	return CMTime{
		CMTimeValue: uint64(currentTime),
		CMTimeScale: c.TimeScale,
		CMTimeFlags: KCMTimeFlags_HasBeenRounded,
		CMTimeEpoch: 0,
	}

}
