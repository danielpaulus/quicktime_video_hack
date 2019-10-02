package messages

import (
	"encoding/binary"
)

/*
All this are just guesses about what this could be from looking at the HexDump.
I guessed that this block could be LPCM (linear pulse code modulation) information although i Dont know.
Seems like 0x 00 00 00 00 00 70 E7 40 6B is some kind of separator or constant, as it is repeated 3 times.
0x6D 63 70 6C renders to "mcpl" (or lpcm in bigendian) in ascii. followed by a few ints.

HexDump:
00 00 00 00 00 70 E7 40 6D 63 70 6C 0C 00 00 00 04 00 00 00 01 00 00 00 04 00 00 00 02 00 00 00 10 00 00 00 00 00 00 00 00 00 00 00 00 70 E7 40 00 00 00 00 00 70 E7 40 6B

*/

const (
	separator uint64 = 0x40E7700000000000
	lpcmMagic uint32 = 0x6C70636D
)

func createLpcmInfo() []byte {
	lpcmBytes := make([]byte, 56)
	binary.LittleEndian.PutUint64(lpcmBytes, separator)
	var index = 8
	binary.LittleEndian.PutUint32(lpcmBytes[index:], lpcmMagic)
	index += 4

	binary.LittleEndian.PutUint32(lpcmBytes[index:], 12)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 4)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 1)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 4)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 2)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 16)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 0)
	index += 4

	binary.LittleEndian.PutUint64(lpcmBytes[index:], separator)
	index += 8
	binary.LittleEndian.PutUint64(lpcmBytes[index:], separator)
	return lpcmBytes
}
