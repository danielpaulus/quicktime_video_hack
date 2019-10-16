package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncOgPacket represents the OG Message. I do not know what these messages mean.
type SyncOgPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
	Unknown       uint32
}

//NewSyncOgPacketFromBytes parses a SyncOgPacket form bytes assuming it starts with SYNC magic and has the correct length.
func NewSyncOgPacketFromBytes(data []byte) (SyncOgPacket, error) {
	packet := SyncOgPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Og Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != OG {
		return packet, fmt.Errorf("wrong message type for OG message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	packet.Unknown = binary.LittleEndian.Uint32(data[24:])
	return packet, nil
}

//NewReply returns a []byte containing the default reply for a SyncOgPacket
func (sp SyncOgPacket) NewReply() []byte {
	responseBytes := make([]byte, 24)
	binary.LittleEndian.PutUint32(responseBytes, 24)
	binary.LittleEndian.PutUint32(responseBytes[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(responseBytes[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint64(responseBytes[16:], 0)

	return responseBytes

}

func (sp SyncOgPacket) String() string {
	return fmt.Sprintf("SYNC_OG{ClockRef:%x, CorrelationID:%x, Unknown:%d}", sp.ClockRef, sp.CorrelationID, sp.Unknown)
}
