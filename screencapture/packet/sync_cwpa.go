package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncCwpaPacket contains all info from a CWPA packet sent by the device
type SyncCwpaPacket struct {
	SyncMagic      uint32
	ClockRef       CFTypeID
	MessageType    uint32
	CorrelationID  uint64
	DeviceClockRef CFTypeID
}

//NewSyncCwpaPacketFromBytes parses a SyncCwpaPacket from a []byte
func NewSyncCwpaPacketFromBytes(data []byte) (SyncCwpaPacket, error) {
	var packet = SyncCwpaPacket{}

	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC CWPA Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	if packet.ClockRef != EmptyCFType {
		return packet, fmt.Errorf("CWPA packet should have empty CFTypeID for ClockRef but has:%x", packet.ClockRef)
	}
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != CWPA {
		return packet, fmt.Errorf("wrong message type for CWPA message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	packet.DeviceClockRef = binary.LittleEndian.Uint64(data[24:])
	return packet, nil
}

//NewReply creates a RPLY packet containing the given clockRef and serializes it to a []byte
func (sp SyncCwpaPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncCwpaPacket) String() string {
	return fmt.Sprintf("SYNC_CWPA{ClockRef:%x, CorrelationID:%x, DeviceClockRef:%x}", sp.ClockRef, sp.CorrelationID, sp.DeviceClockRef)
}
