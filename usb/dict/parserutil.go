package dict

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func writeLengthAndMagic(bytes []byte, length int, magic uint32) {
	binary.LittleEndian.PutUint32(bytes, uint32(length))
	binary.LittleEndian.PutUint32(bytes[4:], magic)
}

func parseLengthAndMagic(bytes []byte, exptectedMagic uint32) (int, []byte, error) {
	length := binary.LittleEndian.Uint32(bytes)
	magic := binary.LittleEndian.Uint32(bytes[4:])
	if int(length) > len(bytes) {
		return 0, bytes, fmt.Errorf("invalid length in header: %d but only received: %d bytes", length, len(bytes))
	}
	if magic != exptectedMagic {
		unknownMagic := string(bytes[4:8])
		return 0, nil, fmt.Errorf("unknown magic type:%s (0x%x), cannot parse value %s", unknownMagic, magic, hex.Dump(bytes))
	}
	return int(length), bytes[8:], nil
}
