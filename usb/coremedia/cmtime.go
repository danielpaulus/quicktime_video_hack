package coremedia

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	KCMTimeFlags_Valid                 uint32 = 0x0
	KCMTimeFlags_HasBeenRounded        uint32 = 0x1
	KCMTimeFlags_PositiveInfinity      uint32 = 0x2
	KCMTimeFlags_NegativeInfinity      uint32 = 0x4
	KCMTimeFlags_Indefinite            uint32 = 0x8
	KCMTimeFlags_ImpliedValueFlagsMask uint32 = KCMTimeFlags_PositiveInfinity | KCMTimeFlags_NegativeInfinity | KCMTimeFlags_Indefinite
	CMTimeLengthInBytes                int    = 24
)

//Taken from https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.8.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMTime.h
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

func (time CMTime) Seconds() uint64 {
	return time.CMTimeValue / uint64(time.CMTimeScale)
}

func (time CMTime) Serialize(target []byte) error {
	if len(target) < CMTimeLengthInBytes {
		return fmt.Errorf("Serializing CMTime failed, not enough space in byte slice:%d", len(target))
	}
	binary.LittleEndian.PutUint64(target, time.CMTimeValue)
	binary.LittleEndian.PutUint32(target[8:], time.CMTimeScale)
	binary.LittleEndian.PutUint32(target[12:], time.CMTimeFlags)
	binary.LittleEndian.PutUint64(target[16:], time.CMTimeEpoch)
	return nil
}

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
	case KCMTimeFlags_Valid:
		flags = "KCMTimeFlags_Valid"
	case KCMTimeFlags_HasBeenRounded:
		flags = "KCMTimeFlags_HasBeenRounded"
	case KCMTimeFlags_PositiveInfinity:
		flags = "KCMTimeFlags_PositiveInfinity"
	case KCMTimeFlags_NegativeInfinity:
		flags = "KCMTimeFlags_NegativeInfinity"
	case KCMTimeFlags_Indefinite:
		flags = "KCMTimeFlags_Indefinite"
	case KCMTimeFlags_ImpliedValueFlagsMask:
		flags = "KCMTimeFlags_ImpliedValueFlagsMask"
	default:
		flags = "unknown"
	}
	return fmt.Sprintf("CMTime{%ds, flags:%s, epoch:%d}", time.Seconds(), flags, time.CMTimeEpoch)
}
