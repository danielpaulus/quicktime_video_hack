package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncClokPacket contains a decoded Clok packet from the device
type SyncClokPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
}

//NewSyncClokPacketFromBytes parses a SynClokPacket from bytes
func NewSyncClokPacketFromBytes(data []byte) (SyncClokPacket, error) {
	packet := SyncClokPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Clok Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != CLOK {
		return packet, fmt.Errorf("wrong message type for Clok message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	return packet, nil
}

//NewReply creates a RPLY message containing the given clockRef and serializes it into a []byte
func (sp SyncClokPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncClokPacket) String() string {
	return fmt.Sprintf("SYNC_CLOK{ClockRef:%x, CorrelationID:%x}", sp.ClockRef, sp.CorrelationID)
}
