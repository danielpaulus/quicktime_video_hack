package packet

import (
	"encoding/binary"
	"fmt"
)

//AsynRelsPacket tells us that a clock was released
type AsynRelsPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
}

//NewAsynRelsPacketFromBytes creates a new AsynRelsPacket from bytes
func NewAsynRelsPacketFromBytes(data []byte) (AsynRelsPacket, error) {
	var packet = AsynRelsPacket{}
	packet.AsyncMagic = binary.LittleEndian.Uint32(data)
	if packet.AsyncMagic != AsynPacketMagic {
		return packet, fmt.Errorf("invalid asyn magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != RELS {
		return packet, fmt.Errorf("invalid packet type in asyn rels:%x", data)
	}

	return packet, nil
}

func (sp AsynRelsPacket) String() string {
	return fmt.Sprintf("ASYN_RELS{ClockRef:%x}", sp.ClockRef)
}
