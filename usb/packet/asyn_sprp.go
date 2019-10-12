package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

type AsynSprpPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	Property    dict.StringKeyEntry
}

func NewAsynSprpPacketFromBytes(data []byte) (AsynSprpPacket, error) {
	var packet = AsynSprpPacket{}
	packet.AsyncMagic = binary.LittleEndian.Uint32(data)
	if packet.AsyncMagic != AsynPacketMagic {
		return packet, fmt.Errorf("invalid asyn magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != SPRP {
		return packet, fmt.Errorf("invalid packet type in asyn sprp:%x", data)
	}
	entry, err := dict.ParseKeyValueEntry(data[16:])
	if err != nil {
		return packet, err
	}
	packet.Property = entry
	return packet, nil
}
