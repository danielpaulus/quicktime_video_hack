package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//SyncTimePacket contains the data from a decoded Time Packet sent by the device
type SyncTimePacket struct {
	ClockRef      CFTypeID
	CorrelationID uint64
}

//NewSyncTimePacketFromBytes parses a SyncTimePacket from bytes
func NewSyncTimePacketFromBytes(data []byte) (SyncTimePacket, error) {
	_, clockRef, correlationID, err := ParseSyncHeader(data, TIME)
	if err != nil {
		return SyncTimePacket{}, err
	}
	packet := SyncTimePacket{ClockRef: clockRef, CorrelationID: correlationID}
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
