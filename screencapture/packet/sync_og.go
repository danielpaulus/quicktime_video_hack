package packet

import (
	"encoding/binary"
	"fmt"
)

//SyncOgPacket represents the OG Message. I do not know what these messages mean.
type SyncOgPacket struct {
	ClockRef      CFTypeID
	CorrelationID uint64
	Unknown       uint32
}

//NewSyncOgPacketFromBytes parses a SyncOgPacket form bytes assuming it starts with SYNC magic and has the correct length.
func NewSyncOgPacketFromBytes(data []byte) (SyncOgPacket, error) {
	remainingBytes, clockRef, correlationID, err := ParseSyncHeader(data, OG)
	if err != nil {
		return SyncOgPacket{}, err
	}
	packet := SyncOgPacket{ClockRef: clockRef, CorrelationID: correlationID}

	packet.Unknown = binary.LittleEndian.Uint32(remainingBytes)
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
