package packet

import (
	"encoding/binary"
)

//Different Sync Packet Magic Markers
const (
	SyncPacketMagic  uint32 = 0x73796E63
	ReplyPacketMagic uint32 = 0x72706C79
	TIME             uint32 = 0x74696D65
	CWPA             uint32 = 0x63777061
	AFMT             uint32 = 0x61666D74
	CVRP             uint32 = 0x63767270
	CLOK             uint32 = 0x636C6F6B
	OG               uint32 = 0x676F2120
)

type CFTypeID = uint64

const EmptyCFType uint64 = 1

func clockRefReply(clockRef uint64, correlationId uint64) []byte {
	length := 28
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data, uint32(length))
	binary.LittleEndian.PutUint32(data[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(data[8:], correlationId)
	binary.LittleEndian.PutUint32(data[16:], 0)
	binary.LittleEndian.PutUint64(data[20:], clockRef)
	return data
}
