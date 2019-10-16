package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/dict"
)

// SyncAfmtPacket contains what I think is information about the audio format
type SyncAfmtPacket struct {
	SyncMagic     uint32
	ClockRef      CFTypeID
	MessageType   uint32
	CorrelationID uint64
	Unknown1      uint32
	Unknown2      uint32
	LpcmMagic     uint32
	LpcmData      coremedia.LPCMData
}

func (sp SyncAfmtPacket) String() string {
	return fmt.Sprintf("SYNC_AFMT{ClockRef:%x, CorrelationID:%x, Unknown1:%x, Unknown2:%x, Lpcm:%s}",
		sp.ClockRef, sp.CorrelationID, sp.Unknown1, sp.Unknown2, sp.LpcmData.String())
}

// NewSyncAfmtPacketFromBytes parses a new AsynFmtPacket from byte array
func NewSyncAfmtPacketFromBytes(data []byte) (SyncAfmtPacket, error) {
	var packet = SyncAfmtPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid sync magic: %x", data)
	}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != AFMT {
		return packet, fmt.Errorf("invalid packet type in sync afmt:%x", data)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	packet.Unknown1 = binary.LittleEndian.Uint32(data[24:])
	packet.Unknown2 = binary.LittleEndian.Uint32(data[28:])
	packet.LpcmMagic = binary.LittleEndian.Uint32(data[32:])
	var err error
	packet.LpcmData, err = coremedia.NewLPCMDataFromBytes(data[36:])
	if err != nil {
		return packet, fmt.Errorf("Error parsing LPCM data in asyn afmt: %s, ", err)
	}
	return packet, nil
}

//NewReply returns a []byte containing a correct reploy for afmt
func (sp SyncAfmtPacket) NewReply() []byte {
	responseDict := createResponseDict()
	dictBytes := dict.SerializeStringKeyDict(responseDict)
	dictLength := uint32(len(dictBytes))
	length := uint32(dictLength + 20)
	responseBytes := make([]byte, length)
	binary.LittleEndian.PutUint32(responseBytes, length)
	binary.LittleEndian.PutUint32(responseBytes[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(responseBytes[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint32(responseBytes[16:], 0)

	copy(responseBytes[20:], dictBytes)
	return responseBytes

}

func createResponseDict() dict.StringKeyDict {
	var response dict.StringKeyDict
	errorCode := common.NewNSNumberFromUInt32(0)
	key := "Error"
	response = dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 1)}
	response.Entries[0].Key = key
	response.Entries[0].Value = errorCode
	return response
}
