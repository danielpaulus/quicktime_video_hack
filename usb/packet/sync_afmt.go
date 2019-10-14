package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
)

//| 4 Byte Length (68)   |4 Byte Magic (SYNC)   | 8 bytes clock CFTypeID| 4 byte magic (AFMT)| 8 byte correlation id
// | some weird data| 4 byte magic (LPCM) |  28 bytes what i think is pcm data|
//AsynAfmtPacket contains what I think is information about the audio format
type SyncAfmtPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationId uint64
	Unknown1      uint32
	Unknown2      uint32
	LpcmMagic     uint32
	LpcmData      messages.LPCMData
}

//NewAsynAfmtPacketFromBytes parses a new AsynFmtPacket from byte array
func NewSyncAfmtPacketFromBytes(data []byte) (SyncAfmtPacket, error) {
	var packet = SyncAfmtPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid sync magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != AFMT {
		return packet, fmt.Errorf("invalid packet type in sync afmt:%x", data)
	}
	packet.CorrelationId = binary.LittleEndian.Uint64(data[16:])
	packet.Unknown1 = binary.LittleEndian.Uint32(data[24:])
	packet.Unknown2 = binary.LittleEndian.Uint32(data[28:])
	packet.LpcmMagic = binary.LittleEndian.Uint32(data[32:])
	var err error
	packet.LpcmData, err = messages.NewLPCMDataFromBytes(data[36:])
	if err != nil {
		return packet, fmt.Errorf("Error parsing LPCM data in asyn afmt: %s, ", err)
	}
	return packet, nil
}
