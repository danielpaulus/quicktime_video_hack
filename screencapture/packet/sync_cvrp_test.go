package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

const expectedCvrpString = "SYNC_CVRP{ClockRef:1, CorrelationID:1135659d0, DeviceClockRef:113538da0, Payload:StringKeyDict:[{PreparedQueueHighWaterLevel : StringKeyDict:[{flags : Int32[1]},{value : UInt64[5]},{timescale : Int32[30]},{epoch : UInt64[0]},]},{PreparedQueueLowWaterLevel : StringKeyDict:[{flags : Int32[1]},{value : UInt64[3]},{timescale : Int32[30]},{epoch : UInt64[0]},]},{FormatDescription : fdsc:{MediaType:Video, VideoDimension:(1126x2436), Codec:AVC-1, PPS:27640033ac5680470133e69e6e04040404, SPS:28ee3cb0, Extensions:IndexKeyDict:[{49 : IndexKeyDict:[{105 : 0x01640033ffe1001127640033ac5680470133e69e6e0404040401000428ee3cb0fdf8f800},]},{52 : H.264},]}},]}"

func TestCvrp(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/cvrp-request")
	if err != nil {
		log.Fatal(err)
	}
	cvrp, err := packet.NewSyncCvrpPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(cvrp.Payload.Entries))
		assert.Equal(t, packet.EmptyCFType, cvrp.ClockRef)
		assert.Equal(t, uint64(0x113538da0), cvrp.DeviceClockRef)
		assert.Equal(t, uint64(0x1135659d0), cvrp.CorrelationID)
		assert.Equal(t, expectedCvrpString, cvrp.String())
	}
	testSerializationOfCvrpReply(cvrp, t)
	_, err = packet.NewSyncCvrpPacketFromBytes(dat)
	assert.Error(t, err)
}

func testSerializationOfCvrpReply(cvrp packet.SyncCvrpPacket, t *testing.T) {
	var clockRef packet.CFTypeID = 0x00007FA66CD10250
	replyBytes := cvrp.NewReply(clockRef)
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/cvrp-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}
