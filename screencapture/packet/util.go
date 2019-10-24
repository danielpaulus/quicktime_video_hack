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

//ParseSyncHeader checks for the SYNC magic and the given messagemagic and then returns
//the remainingBytes starting after the messagemagic so after 16 bytes, the clockRef, correlationID and an error
//is the packet is not SYNC or if the messagemagic is wrong
func ParseSyncHeader(data []byte, messagemagic uint32) ([]byte, CFTypeID, uint64, error) {
	remainingBytes, clockRef, err := parseHeader(data, SyncPacketMagic, messagemagic)
	if err != nil {
		return data, 0, 0, err
	}
	correlationID := binary.LittleEndian.Uint64(remainingBytes)
	return remainingBytes[8:], clockRef, correlationID, err
}

func parseHeader(data []byte, packetmagic uint32, messagemagic uint32) ([]byte, CFTypeID, error) {
	magic := binary.LittleEndian.Uint32(data)
	if magic != packetmagic {
		packetTypeASCII := string(data[:4])
		return nil, 0, fmt.Errorf("invalid packet magic '%s' - packethex: %x", packetTypeASCII, data)
	}
	clockRef := binary.LittleEndian.Uint64(data[4:])
	messageType := binary.LittleEndian.Uint32(data[12:])
	if messageType != messagemagic {
		messageTypeASCII := string(data[12:16])
		return nil, 0, fmt.Errorf("invalid packet type:'%s' - packethex: %x", messageTypeASCII, data)
	}
	return data[16:], clockRef, nil
}
