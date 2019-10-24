package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncStopPacket requests us to stop our clock
type SyncStopPacket struct {
	ClockRef      CFTypeID
	CorrelationID uint64
}

//NewSyncStopPacketFromBytes parses a SyncStopPacket from bytes
func NewSyncStopPacketFromBytes(data []byte) (SyncStopPacket, error) {
	_, clockRef, correlationID, err := ParseSyncHeader(data, STOP)
	if err != nil {
		return SyncStopPacket{}, err
	}
	packet := SyncStopPacket{ClockRef: clockRef, CorrelationID: correlationID}
	return packet, nil
}

//NewReply creates a byte array containing the given skew
func (sp SyncStopPacket) NewReply() []byte {
	responseBytes := make([]byte, 24)
	binary.LittleEndian.PutUint32(responseBytes, 24)
	binary.LittleEndian.PutUint32(responseBytes[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(responseBytes[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint32(responseBytes[16:], 0)
	binary.LittleEndian.PutUint32(responseBytes[20:], 0)
	return responseBytes
}

func (sp SyncStopPacket) String() string {
	return fmt.Sprintf("SYNC_STOP{ClockRef:%x, CorrelationID:%x}", sp.ClockRef, sp.CorrelationID)
}
