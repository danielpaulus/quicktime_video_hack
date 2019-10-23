package packet

import (
	"encoding/binary"
	"fmt"
)

//AsynTbasPacket contains info about a new Timebase. I do not know what the other reference is used for.
type AsynTbasPacket struct {
	ClockRef     CFTypeID
	SomeOtherRef CFTypeID
}

//NewAsynTbasPacketFromBytes parses a AsynTbasPacket from bytes.
func NewAsynTbasPacketFromBytes(data []byte) (AsynTbasPacket, error) {
	var packet = AsynTbasPacket{}
	remainingBytes, clockRef, err := ParseAsynHeader(data, TBAS)
	if err != nil {
		return packet, err
	}
	packet.ClockRef = clockRef
	packet.SomeOtherRef = binary.LittleEndian.Uint64(remainingBytes)
	return packet, nil
}

func (sp AsynTbasPacket) String() string {
	return fmt.Sprintf("ASYN_TBAS{ClockRef:%x, UnknownRef:%x}", sp.ClockRef, sp.SomeOtherRef)
}
