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

//CFTypeID is just a type alias for uint64 but I think it is closer to what is happening on MAC/iOS
type CFTypeID = uint64

//EmptyCFType is a CFTypeId of 0x1
const EmptyCFType CFTypeID = 1

func clockRefReply(clockRef uint64, correlationID uint64) []byte {
	length := 28
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data, uint32(length))
	binary.LittleEndian.PutUint32(data[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(data[8:], correlationID)
	binary.LittleEndian.PutUint32(data[16:], 0)
	binary.LittleEndian.PutUint64(data[20:], clockRef)
	return data
}
