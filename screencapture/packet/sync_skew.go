package packet

import (
	"encoding/binary"
	"fmt"
	"math"
)

//SyncSkewPacket requests us to reply with the current skew value
type SyncSkewPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
}

//NewSyncSkewPacketFromBytes parses a SyncSkewPacket from bytes
func NewSyncSkewPacketFromBytes(data []byte) (SyncSkewPacket, error) {
	packet := SyncSkewPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Skew Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != SKEW {
		return packet, fmt.Errorf("wrong message type for SKEW message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
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
