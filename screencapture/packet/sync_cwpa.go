package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncCwpaPacket contains all info from a CWPA packet sent by the device
type SyncCwpaPacket struct {
	ClockRef       CFTypeID
	CorrelationID  uint64
	DeviceClockRef CFTypeID
}

//NewSyncCwpaPacketFromBytes parses a SyncCwpaPacket from a []byte
func NewSyncCwpaPacketFromBytes(data []byte) (SyncCwpaPacket, error) {
	remainingBytes, clockRef, correlationID, err := ParseSyncHeader(data, CWPA)
	if err != nil {
		return SyncCwpaPacket{}, err
	}
	packet := SyncCwpaPacket{ClockRef: clockRef, CorrelationID: correlationID}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	if packet.ClockRef != EmptyCFType {
		return packet, fmt.Errorf("CWPA packet should have empty CFTypeID for ClockRef but has:%x", packet.ClockRef)
	}

	packet.DeviceClockRef = binary.LittleEndian.Uint64(remainingBytes)
	return packet, nil
}

//NewReply creates a RPLY packet containing the given clockRef and serializes it to a []byte
func (sp SyncCwpaPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncCwpaPacket) String() string {
	return fmt.Sprintf("SYNC_CWPA{ClockRef:%x, CorrelationID:%x, DeviceClockRef:%x}", sp.ClockRef, sp.CorrelationID, sp.DeviceClockRef)
}
