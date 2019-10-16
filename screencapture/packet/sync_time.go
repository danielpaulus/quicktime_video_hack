package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//SyncTimePacket contains the data from a decoded Time Packet sent by the device
type SyncTimePacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
}

//NewSyncTimePacketFromBytes parses a SyncTimePacket from bytes
func NewSyncTimePacketFromBytes(data []byte) (SyncTimePacket, error) {
	packet := SyncTimePacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Time Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != TIME {
		return packet, fmt.Errorf("wrong message type for Time message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	return packet, nil
}

//NewReply creates a RPLY packet containing the given CMTime and serializes it to a []byte
func (sp SyncTimePacket) NewReply(time coremedia.CMTime) ([]byte, error) {
	length := 44
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data, uint32(length))
	binary.LittleEndian.PutUint32(data[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(data[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint32(data[16:], 0)
	err := time.Serialize(data[20:])
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (sp SyncTimePacket) String() string {
	return fmt.Sprintf("SYNC_TIME{ClockRef:%x, CorrelationID:%x}", sp.ClockRef, sp.CorrelationID)
}
