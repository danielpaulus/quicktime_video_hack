package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

//AsynSprpPacket seems to be a set property packet sent by the device.
type AsynSprpPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	Property    dict.StringKeyEntry
}

//NewAsynSprpPacketFromBytes creates a new AsynSprpPacket from bytes
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

func (sp AsynSprpPacket) String() string {
	return fmt.Sprintf("ASYN_SPRP{ClockRef:%x, Property:{%s:%s}}", sp.ClockRef, sp.Property.Key, sp.Property.Value)
}
