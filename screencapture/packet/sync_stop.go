package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncStopPacket requests us to stop our clock
type SyncStopPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
}

//NewSyncStopPacketFromBytes parses a SyncStopPacket from bytes
func NewSyncStopPacketFromBytes(data []byte) (SyncStopPacket, error) {
	packet := SyncStopPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC STOP Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != STOP {
		return packet, fmt.Errorf("wrong message type for STOP message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
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
