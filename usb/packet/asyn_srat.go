package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/usb/coremedia"
	"math"
)

// AsynSratPacket is probably related to AVPlayer.SetRate somehow. I dont know exactly what everything means here
type AsynSratPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	Rate1       float32
	Rate2       float32
	Time        coremedia.CMTime
}

//NewAsynSratPacketFromBytes parses a new AsynSratPacket from bytes
func NewAsynSratPacketFromBytes(data []byte) (AsynSratPacket, error) {
	var packet = AsynSratPacket{}
	packet.AsyncMagic = binary.LittleEndian.Uint32(data)
	if packet.AsyncMagic != AsynPacketMagic {
		return packet, fmt.Errorf("invalid asyn magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != SRAT {
		return packet, fmt.Errorf("invalid packet type in asyn tjmp:%x", data)
	}

	packet.Rate1 = math.Float32frombits(binary.LittleEndian.Uint32(data[16:]))
	packet.Rate2 = math.Float32frombits(binary.LittleEndian.Uint32(data[20:]))
	cmtime, err := coremedia.NewCMTimeFromBytes(data[24:])
	if err != nil {
		return packet, err
	}
	packet.Time = cmtime
	return packet, nil
}

func (sp AsynSratPacket) String() string {
	return fmt.Sprintf("ASYN_TBAS{ClockRef:%x, Rate1:%f, Rate2:%f, Time:%s}", sp.ClockRef, sp.Rate1, sp.Rate2, sp.Time.String())
}
