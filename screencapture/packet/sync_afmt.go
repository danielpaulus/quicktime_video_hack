package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

// SyncAfmtPacket contains what I think is information about the audio format
type SyncAfmtPacket struct {
	SyncMagic                   uint32
	ClockRef                    CFTypeID
	MessageType                 uint32
	CorrelationID               uint64
	AudioStreamBasicDescription coremedia.AudioStreamBasicDescription
}

func (sp SyncAfmtPacket) String() string {
	return fmt.Sprintf("SYNC_AFMT{ClockRef:%x, CorrelationID:%x, AudioStreamBasicDescription:%s}",
		sp.ClockRef, sp.CorrelationID, sp.AudioStreamBasicDescription.String())
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

	var err error
	packet.AudioStreamBasicDescription, err = coremedia.NewAudioStreamBasicDescriptionFromBytes(data[24:])
	if err != nil {
		return packet, fmt.Errorf("Error parsing AudioStreamBasicDescription data in asyn afmt: %s, ", err)
	}
	return packet, nil
}

//NewReply returns a []byte containing a correct reploy for afmt
func (sp SyncAfmtPacket) NewReply() []byte {
	responseDict := createResponseDict()
	dictBytes := coremedia.SerializeStringKeyDict(responseDict)
	dictLength := uint32(len(dictBytes))
	length := dictLength + 20
	responseBytes := make([]byte, length)
	binary.LittleEndian.PutUint32(responseBytes, length)
	binary.LittleEndian.PutUint32(responseBytes[4:], ReplyPacketMagic)
	binary.LittleEndian.PutUint64(responseBytes[8:], sp.CorrelationID)
	binary.LittleEndian.PutUint32(responseBytes[16:], 0)

	copy(responseBytes[20:], dictBytes)
	return responseBytes

}

func createResponseDict() coremedia.StringKeyDict {
	var response coremedia.StringKeyDict
	errorCode := common.NewNSNumberFromUInt32(0)
	key := "Error"
	response = coremedia.StringKeyDict{Entries: make([]coremedia.StringKeyEntry, 1)}
	response.Entries[0].Key = key
	response.Entries[0].Value = errorCode
	return response
}
