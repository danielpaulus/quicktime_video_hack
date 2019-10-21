package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//SyncCvrpPacket contains all info from a CVRP packet sent by the device
type SyncCvrpPacket struct {
	SyncMagic      uint32
	ClockRef       CFTypeID
	MessageType    uint32
	CorrelationID  uint64
	DeviceClockRef CFTypeID
	Payload        coremedia.StringKeyDict
}

//NewSyncCvrpPacketFromBytes parses a SyncCvrpPacket from a []byte
func NewSyncCvrpPacketFromBytes(data []byte) (SyncCvrpPacket, error) {
	packet := SyncCvrpPacket{}
	packet.SyncMagic = binary.LittleEndian.Uint32(data)
	if packet.SyncMagic != SyncPacketMagic {
		return packet, fmt.Errorf("invalid SYNC Cvrp Packet: %x", data)
	}

	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	if packet.ClockRef != EmptyCFType {
		return packet, fmt.Errorf("Cvrp packet should have empty CFTypeID for ClockRef but has:%x", packet.ClockRef)
	}
	packet.MessageType = binary.LittleEndian.Uint32(data[12:])
	if packet.MessageType != CVRP {
		return packet, fmt.Errorf("wrong message type for Cvrp message: %x", packet.MessageType)
	}
	packet.CorrelationID = binary.LittleEndian.Uint64(data[16:])
	packet.DeviceClockRef = binary.LittleEndian.Uint64(data[24:])
	payloadDict, err := coremedia.NewStringDictFromBytes(data[32:])
	if err != nil {
		return packet, err
	}
	packet.Payload = payloadDict
	return packet, nil
}

//NewReply creates a RPLY packet containing the given clockRef and serializes it to a []byte
func (sp SyncCvrpPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncCvrpPacket) String() string {
	return fmt.Sprintf("SYNC_CVRP{ClockRef:%x, CorrelationID:%x, DeviceClockRef:%x, Payload:%s}", sp.ClockRef, sp.CorrelationID, sp.DeviceClockRef, sp.Payload.String())
}
