package packet

import (
	"encoding/binary"
	"fmt"
)

//ParseAsynHeader checks for the ASYN magic and the given messagemagic and then returns
//the remainingBytes starting after the messagemagic so after 16 bytes, the clockRef and an error
//is the packet is not Asyn or if the messagemagic is wrong
func ParseAsynHeader(data []byte, messagemagic uint32) ([]byte, CFTypeID, error) {
	return parseHeader(data, AsynPacketMagic, messagemagic)
}

func parseHeader(data []byte, packetmagic uint32, messagemagic uint32) ([]byte, CFTypeID, error) {
	magic := binary.LittleEndian.Uint32(data)
	if magic != packetmagic {
		return nil, 0, fmt.Errorf("invalid asyn magic: %x", data)
	}
	clockRef := binary.LittleEndian.Uint64(data[4:])
	messageType := binary.LittleEndian.Uint32(data[12:])
	if messageType != messagemagic {
		return nil, 0, fmt.Errorf("invalid packet type in asyn sprp:%x", data)
	}
	return data[16:], clockRef, nil
}
