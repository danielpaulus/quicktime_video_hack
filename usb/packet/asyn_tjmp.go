package packet

import (
	"encoding/binary"
	"fmt"
)

//AsynTjmpPacket contains the data from a TJMP packet.
//I think this is a notification sent by the device about changing a TimeBase.
//I do not know what the last bytes are for currently.
type AsynTjmpPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	Unknown     []byte
}

//NewAsynTjmpPacketFromBytes parses a new AsynTjmpPacket from byte array
func NewAsynTjmpPacketFromBytes(data []byte) (AsynTjmpPacket, error) {
	var packet = AsynTjmpPacket{}
	packet.AsyncMagic = binary.LittleEndian.Uint32(data)
	if packet.AsyncMagic != AsynPacketMagic {
		return packet, fmt.Errorf("invalid asyn magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != TJMP {
		return packet, fmt.Errorf("invalid packet type in asyn tjmp:%x", data)
	}

	packet.Unknown = data[16:]
	return packet, nil
}

func (sp AsynTjmpPacket) String() string {
	return fmt.Sprintf("ASYN_TJMP{ClockRef:%x, UnknownData:%x}", sp.ClockRef, sp.Unknown)
}
