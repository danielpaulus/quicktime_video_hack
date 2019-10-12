package packet

import (
	"encoding/binary"
	"fmt"
)

type AsynTbasPacket struct {
	AsyncMagic   uint32
	ClockRef     CFTypeID
	MessageType  uint32
	SomeOtherRef CFTypeID
}

func NewAsynTbasPacketFromBytes(data []byte) (AsynTbasPacket, error) {
	var packet = AsynTbasPacket{}
	packet.AsyncMagic = binary.LittleEndian.Uint32(data)
	if packet.AsyncMagic != AsynPacketMagic {
		return packet, fmt.Errorf("invalid asyn magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != TBAS {
		return packet, fmt.Errorf("invalid packet type in asyn tbas:%x", data)
	}
	packet.SomeOtherRef = binary.LittleEndian.Uint64(data[16:])
	return packet, nil
}
