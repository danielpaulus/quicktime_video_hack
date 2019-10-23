package packet

import (
	"fmt"
)

//AsynRelsPacket tells us that a clock was released
type AsynRelsPacket struct {
	ClockRef CFTypeID
}

//NewAsynRelsPacketFromBytes creates a new AsynRelsPacket from bytes
func NewAsynRelsPacketFromBytes(data []byte) (AsynRelsPacket, error) {
	var packet = AsynRelsPacket{}
	_, clockRef, err := ParseAsynHeader(data, RELS)
	if err != nil {
		return packet, err
	}
	packet.ClockRef = clockRef

	return packet, nil
}

func (sp AsynRelsPacket) String() string {
	return fmt.Sprintf("ASYN_RELS{ClockRef:%x}", sp.ClockRef)
}
