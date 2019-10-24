package packet

import (
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//AsynSprpPacket seems to be a set property packet sent by the device.
type AsynSprpPacket struct {
	ClockRef CFTypeID
	Property coremedia.StringKeyEntry
}

//NewAsynSprpPacketFromBytes creates a new AsynSprpPacket from bytes
func NewAsynSprpPacketFromBytes(data []byte) (AsynSprpPacket, error) {
	var packet = AsynSprpPacket{}
	remainingBytes, clockRef, err := ParseAsynHeader(data, SPRP)
	if err != nil {
		return packet, err
	}
	packet.ClockRef = clockRef
	entry, err := coremedia.ParseKeyValueEntry(remainingBytes)
	if err != nil {
		return packet, err
	}
	packet.Property = entry
	return packet, nil
}

func (sp AsynSprpPacket) String() string {
	return fmt.Sprintf("ASYN_SPRP{ClockRef:%x, Property:{%s:%s}}", sp.ClockRef, sp.Property.Key, sp.Property.Value)
}
