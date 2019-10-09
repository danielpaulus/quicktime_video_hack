package packet

import (
	"encoding/binary"
	"fmt"
)

type SyncTimePacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
}

func NewSyncTimePacketFromBytes(data []byte) (SyncTimePacket, error) {
	packet := SyncTimePacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Time Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != TIME {
		return packet, fmt.Errorf("wrong message type for Time message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	return packet, nil
}

func (sp SyncTimePacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}
