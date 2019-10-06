package packet

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

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
	Header  uint64 //I don't know what the first 8 bytes are for currently
	Magic   uint32
	Payload dict.StringKeyDict
}

func ExtractDictFromBytes(data []byte) (SyncPacket, error) {
	result := SyncPacket{}
	result.Header = binary.LittleEndian.Uint64(data)
	result.Magic = binary.LittleEndian.Uint32(data[8:])
	payloadDict, err := dict.NewStringDictFromBytes(data[28:])
	if err != nil {
		return result, err
	}
	result.Payload = payloadDict
	return result, nil
}
