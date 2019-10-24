package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

const expectedString = `ASYN_SBUF{ClockRef:7ffb5cc32f60, sBuf:{OutputPresentationTS:CMTime{95911997690984/1000000000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, NumSamples:1, Nalus:[{len:30 type:SEI},{len:90712 type:IDR},], fdsc:fdsc:{MediaType:Video, VideoDimension:(1126x2436), Codec:AVC-1, PPS:27640033ac5680470133e69e6e04040404, SPS:28ee3cb0, Extensions:IndexKeyDict:[{49 : IndexKeyDict:[{105 : 0x01640033ffe1001127640033ac5680470133e69e6e0404040401000428ee3cb0fdf8f800},]},{52 : H.264},]}, attach:IndexKeyDict:[{28 : IndexKeyDict:[{46 : Float64[2436.000000]},{47 : Float64[2436.000000]},]},{29 : Int32[0]},{26 : IndexKeyDict:[{46 : Float64[1126.000000]},{47 : Float64[2436.000000]},{45 : Float64[0.000000]},{44 : Float64[0.000000]},]},{27 : IndexKeyDict:[{46 : Float64[1126.000000]},{47 : Float64[2436.000000]},{45 : Float64[0.000000]},{44 : Float64[0.000000]},]},], sary:IndexKeyDict:[{4 : %!s(bool=false)},], SampleTimingInfoArray:{Duration:CMTime{1/60, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, PresentationTS:CMTime{95911997690984/1000000000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, DecodeTS:CMTime{0/0, flags:KCMTimeFlagsValid, epoch:0}}}}`

func TestFeed(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}
	feedPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x7ffb5cc32f60), feedPacket.ClockRef)
		assert.Equal(t, expectedString, feedPacket.String())
		assert.Equal(t, coremedia.MediaTypeVideo, feedPacket.CMSampleBuf.MediaType)
	}
}

func TestEat(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/asyn-eat")
	if err != nil {
		log.Fatal(err)
	}
	feedPacket, err := packet.NewAsynCmSampleBufPacketFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0x133959728), feedPacket.ClockRef)
		assert.Equal(t, coremedia.MediaTypeSound, feedPacket.CMSampleBuf.MediaType)
	}
}
