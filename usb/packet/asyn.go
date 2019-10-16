package packet

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

//Async Packet types
const (
	AsynPacketMagic uint32 = 0x6173796E
	FEED            uint32 = 0x66656564 //These contain CMSampleBufs which contain raw h264 Nalus
	TJMP            uint32 = 0x746A6D70
	SRAT            uint32 = 0x73726174 //CMTimebaseSetRateAndAnchorTime https://developer.apple.com/documentation/coremedia/cmtimebase?language=objc
	SPRP            uint32 = 0x73707270 // Set Property
	TBAS            uint32 = 0x74626173 //TimeBase https://developer.apple.com/library/archive/qa/qa1643/_index.html
	RELS            uint32 = 0x72656C73
	HPD1            uint32 = 0x68706431 //hpd1 - 1dph | For specifying/requesting the video format
	HPA1            uint32 = 0x68706131 //hpa1 - 1aph | For specifying/requesting the audio format
	NEED            uint32 = 0x6E656564 //need - deen
	EAT             uint32 = 0x65617421 //contains audio sbufs
)

type AsyncPacket struct {
	Header                     uint64 //I don't know what the first 8 bytes are for currently
	HumanReadableTypeSpecifier uint32 //One of the packet types above
	Payload                    interface{}
}

func NewAsynHpd1Packet(stringKeyDict dict.StringKeyDict) []byte {
	return newAsynDictPacket(stringKeyDict, HPD1, EmptyCFType)
}

func NewAsynHpa1Packet(stringKeyDict dict.StringKeyDict, clockRef CFTypeID) []byte {
	return newAsynDictPacket(stringKeyDict, HPA1, clockRef)
}

func newAsynDictPacket(stringKeyDict dict.StringKeyDict, subtypeMarker uint32, asynTypeHeader uint64) []byte {
	serialize := dict.SerializeStringKeyDict(stringKeyDict)
	length := len(serialize) + 20
	header := make([]byte, 20)
	binary.LittleEndian.PutUint32(header, uint32(length))
	binary.LittleEndian.PutUint32(header[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(header[8:], asynTypeHeader)
	binary.LittleEndian.PutUint32(header[16:], subtypeMarker)
	return append(header, serialize...)
}

//AsynNeedPacketBytes can be used to create the NEED message as soon as the clockRef from SYNC CVRP has been received.
func AsynNeedPacketBytes(clockRef CFTypeID) []byte {
	needPacketLength := 20
	packet := make([]byte, needPacketLength)
	binary.LittleEndian.PutUint32(packet, uint32(needPacketLength))
	binary.LittleEndian.PutUint32(packet[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(packet[8:], clockRef)
	binary.LittleEndian.PutUint32(packet[16:], NEED) //need - deen
	return packet
}
