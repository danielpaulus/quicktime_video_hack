package packet

import (
	"fmt"
)

//AsynTjmpPacket contains the data from a TJMP packet.
//I think this is a notification sent by the device about changing a TimeBase.
//I do not know what the last bytes are for currently.
type AsynTjmpPacket struct {
	ClockRef CFTypeID
	Unknown  []byte
}

//NewAsynTjmpPacketFromBytes parses a new AsynTjmpPacket from byte array
func NewAsynTjmpPacketFromBytes(data []byte) (AsynTjmpPacket, error) {
	var packet = AsynTjmpPacket{}
	remainingBytes, clockRef, err := ParseAsynHeader(data, TJMP)
	if err != nil {
		return packet, err
	}
	packet.ClockRef = clockRef

	packet.Unknown = remainingBytes
	return packet, nil
}

func (sp AsynTjmpPacket) String() string {
	return fmt.Sprintf("ASYN_TJMP{ClockRef:%x, UnknownData:%x}", sp.ClockRef, sp.Unknown)
}
