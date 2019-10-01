package packet

//Different Sync Packet Magic Markers
const (
	SyncPacketMagic uint32 = 0x73796E63
	TIME            uint32 = 0x74696D65
	CWPA            uint32 = 0x63777061
	AFMT            uint32 = 0x61666D74
	CVRP            uint32 = 0x63767270
	CLOK            uint32 = 0x636C6F6B
)

type SyncPacket struct {
	Header                     uint64 //I don't know what the first 8 bytes are for currently
	HumanReadableTypeSpecifier uint32 //One of the packet types above
	Payload                    interface{}
}
