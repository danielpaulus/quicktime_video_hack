package packet

import (
	"fmt"
)

//SyncClokPacket contains a decoded Clok packet from the device
type SyncClokPacket struct {
	ClockRef      CFTypeID
	CorrelationID uint64
}

//NewSyncClokPacketFromBytes parses a SynClokPacket from bytes
func NewSyncClokPacketFromBytes(data []byte) (SyncClokPacket, error) {
	_, clockRef, correlationID, err := ParseSyncHeader(data, CLOK)
	if err != nil {
		return SyncClokPacket{}, err
	}
	packet := SyncClokPacket{ClockRef: clockRef, CorrelationID: correlationID}
	return packet, nil
}

//NewReply creates a RPLY message containing the given clockRef and serializes it into a []byte
func (sp SyncClokPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncClokPacket) String() string {
	return fmt.Sprintf("SYNC_CLOK{ClockRef:%x, CorrelationID:%x}", sp.ClockRef, sp.CorrelationID)
}
