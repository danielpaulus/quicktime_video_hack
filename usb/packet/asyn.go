package packet

//Async Packet types
const (
	AsynPacketMagic uint32 = 0x6173796E
	FEED            uint32 = 0x66656564 //These contain CMSampleBufs which contain raw h264 Nalus
	TJMP            uint32 = 0x746A6D70
	SRAT            uint32 = 0x73726174
	SPRP            uint32 = 0x73707270
	TBAS            uint32 = 0x74626173
	RELS            uint32 = 0x72656C73
)

type AsyncPacket struct {
	Header                     uint64 //I don't know what the first 8 bytes are for currently
	HumanReadableTypeSpecifier uint32 //One of the packet types above
	Payload                    interface{}
}
