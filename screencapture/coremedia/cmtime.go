package coremedia

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//Constants for the CMTime struct
const (
	KCMTimeFlagsValid                 uint32 = 0x0
	KCMTimeFlagsHasBeenRounded        uint32 = 0x1
	KCMTimeFlagsPositiveInfinity      uint32 = 0x2
	KCMTimeFlagsNegativeInfinity      uint32 = 0x4
	KCMTimeFlagsIndefinite            uint32 = 0x8
	KCMTimeFlagsImpliedValueFlagsMask        = KCMTimeFlagsPositiveInfinity | KCMTimeFlagsNegativeInfinity | KCMTimeFlagsIndefinite
	CMTimeLengthInBytes               int    = 24
)

//CMTime is taken from https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.8.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMTime.h
type CMTime struct {
	CMTimeValue uint64 /*! @field value The value of the CMTime. value/timescale = seconds. */
	CMTimeScale uint32 /*! @field timescale The timescale of the CMTime. value/timescale = seconds.  */
	CMTimeFlags uint32 /*! @field flags The flags, eg. kCMTimeFlags_Valid, kCMTimeFlags_PositiveInfinity, etc. */
	CMTimeEpoch uint64 /*! @field epoch Differentiates between equal timestamps that are actually different because
	of looping, multi-item sequencing, etc.
	Will be used during comparison: greater epochs happen after lesser ones.
	Additions/subtraction is only possible within a single epoch,
	however, since epoch length may be unknown/variable. */
}

//GetTimeForScale calculates a float64 TimeValue by rescaling this CMTime to the CMTimeScale of the given CMTime
func (time CMTime) GetTimeForScale(newScaleToUse CMTime) float64 {
	scalingFactor := float64(newScaleToUse.CMTimeScale) / float64(time.CMTimeScale)
	return (float64(time.CMTimeValue) * scalingFactor)
}

//Seconds returns CMTimeValue/CMTimeScale and 0 when all values are 0
func (time CMTime) Seconds() uint64 {
	//prevent division by 0
	if time.CMTimeValue == 0 {
		return 0
	}
	return time.CMTimeValue / uint64(time.CMTimeScale)
}

//Serialize serializes this CMTime into a given byte slice that needs to be at least of CMTimeLengthInBytes length
func (time CMTime) Serialize(target []byte) error {
	if len(target) < CMTimeLengthInBytes {
		return fmt.Errorf("serializing CMTime failed, not enough space in byte slice:%d", len(target))
	}
	binary.LittleEndian.PutUint64(target, time.CMTimeValue)
	binary.LittleEndian.PutUint32(target[8:], time.CMTimeScale)
	binary.LittleEndian.PutUint32(target[12:], time.CMTimeFlags)
	binary.LittleEndian.PutUint64(target[16:], time.CMTimeEpoch)
	return nil
}

//NewCMTimeFromBytes reads a CMTime struct directly from the given byte slice
func NewCMTimeFromBytes(data []byte) (CMTime, error) {
	r := bytes.NewReader(data)
	var cmTime CMTime
	err := binary.Read(r, binary.LittleEndian, &cmTime)
	if err != nil {
		return cmTime, err
	}
	return cmTime, nil
}

func (time CMTime) String() string {
	var flags string
	switch time.CMTimeFlags {
	case KCMTimeFlagsValid:
		flags = "KCMTimeFlagsValid"
	case KCMTimeFlagsHasBeenRounded:
		flags = "KCMTimeFlagsHasBeenRounded"
	case KCMTimeFlagsPositiveInfinity:
		flags = "KCMTimeFlagsPositiveInfinity"
	case KCMTimeFlagsNegativeInfinity:
		flags = "KCMTimeFlagsNegativeInfinity"
	case KCMTimeFlagsIndefinite:
		flags = "KCMTimeFlagsIndefinite"
	case KCMTimeFlagsImpliedValueFlagsMask:
		flags = "KCMTimeFlagsImpliedValueFlagsMask"
	default:
		flags = "unknown"
	}
	return fmt.Sprintf("CMTime{%d/%d, flags:%s, epoch:%d}", time.CMTimeValue, time.CMTimeScale, flags, time.CMTimeEpoch)
}
