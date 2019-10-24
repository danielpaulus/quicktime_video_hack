package packet

import (
	"encoding/binary"
	"fmt"
	"math"
)

//SyncSkewPacket requests us to reply with the current skew value
type SyncSkewPacket struct {
	ClockRef      CFTypeID
	CorrelationID uint64
}

//NewSyncSkewPacketFromBytes parses a SyncSkewPacket from bytes
func NewSyncSkewPacketFromBytes(data []byte) (SyncSkewPacket, error) {
	_, clockRef, correlationID, err := ParseSyncHeader(data, SKEW)
	if err != nil {
		return SyncSkewPacket{}, err
	}
	packet := SyncSkewPacket{ClockRef: clockRef, CorrelationID: correlationID}
	return packet, nil
}

//NewReply creates a byte array containing the given skew
func (sp SyncSkewPacket) NewReply(skew float64) []byte {
	responseBytes := make([]byte, 28)
	binary.LittleEndian.PutUint32(responseBytes, 28)
	binary.LittleEndian.PutUint32(responseBytes[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(responseBytes[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint32(responseBytes[16:], 0)
	binary.LittleEndian.PutUint64(responseBytes[20:], math.Float64bits(skew))
	return responseBytes
}

func (sp SyncSkewPacket) String() string {
	return fmt.Sprintf("SYNC_SKEW{ClockRef:%x, CorrelationID:%x}", sp.ClockRef, sp.CorrelationID)
}
