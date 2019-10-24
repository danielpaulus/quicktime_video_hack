package packet

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

// AsynSratPacket is probably related to AVPlayer.SetRate somehow. I dont know exactly what everything means here
type AsynSratPacket struct {
	ClockRef CFTypeID
	Rate1    float32
	Rate2    float32
	Time     coremedia.CMTime
}

//NewAsynSratPacketFromBytes parses a new AsynSratPacket from bytes
func NewAsynSratPacketFromBytes(data []byte) (AsynSratPacket, error) {
	var packet = AsynSratPacket{}
	remainingBytes, clockRef, err := ParseAsynHeader(data, SRAT)
	if err != nil {
		return packet, err
	}
	packet.ClockRef = clockRef

	packet.Rate1 = math.Float32frombits(binary.LittleEndian.Uint32(remainingBytes))
	packet.Rate2 = math.Float32frombits(binary.LittleEndian.Uint32(remainingBytes[4:]))
	cmtime, err := coremedia.NewCMTimeFromBytes(remainingBytes[8:])
	if err != nil {
		return packet, err
	}
	packet.Time = cmtime
	return packet, nil
}

func (sp AsynSratPacket) String() string {
	return fmt.Sprintf("ASYN_SRAT{ClockRef:%x, Rate1:%f, Rate2:%f, Time:%s}", sp.ClockRef, sp.Rate1, sp.Rate2, sp.Time.String())
}
