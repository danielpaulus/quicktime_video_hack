package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

//WriteLengthAndMagic just writes length and magic as uint32 4 byte values into the given array.
func WriteLengthAndMagic(bytes []byte, length int, magic uint32) {
	binary.LittleEndian.PutUint32(bytes, uint32(length))
	binary.LittleEndian.PutUint32(bytes[4:], magic)
}

//ParseLengthAndMagic checks if if the given byte array is longer or equal the uint32 in the first 4 bytes, and if the magic value in the second 4 bytes equals the supplied magic
// and returns the length, a slice of the bytes without length and magic or an error.
func ParseLengthAndMagic(bytes []byte, exptectedMagic uint32) (int, []byte, error) {
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
