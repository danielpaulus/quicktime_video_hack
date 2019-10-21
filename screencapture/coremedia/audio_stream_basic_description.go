package coremedia

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

//AudioFormatIDLpcm is the CoreMedia MediaID for LPCM
const AudioFormatIDLpcm uint32 = 0x6C70636D

//AudioStreamBasicDescription represents the struct found here: https://github.com/nu774/MSResampler/blob/master/CoreAudio/CoreAudioTypes.h
type AudioStreamBasicDescription struct {
	SampleRate       float64
	FormatID         uint32
	FormatFlags      uint32
	BytesPerPacket   uint32
	FramesPerPacket  uint32
	BytesPerFrame    uint32
	ChannelsPerFrame uint32
	BitsPerChannel   uint32
	Reserved         uint32
}

//DefaultAudioStreamBasicDescription creates a LPCM AudioStreamBasicDescription with default values I grabbed from the hex dump
func DefaultAudioStreamBasicDescription() AudioStreamBasicDescription {
	return AudioStreamBasicDescription{FormatFlags: 12,
		BytesPerPacket: 4, FramesPerPacket: 1, BytesPerFrame: 4, ChannelsPerFrame: 2, BitsPerChannel: 16, Reserved: 0,
		SampleRate: 48000, FormatID: AudioFormatIDLpcm}
}

func (adsb AudioStreamBasicDescription) String() string {
	return fmt.Sprintf("{SampleRate:%f,FormatFlags:%d,BytesPerPacket:%d,FramesPerPacket:%d,BytesPerFrame:%d,ChannelsPerFrame:%d,BitsPerChannel:%d,Reserved:%d}",
		adsb.SampleRate, adsb.FormatFlags, adsb.BytesPerPacket, adsb.FramesPerPacket,
		adsb.BytesPerFrame, adsb.ChannelsPerFrame, adsb.BitsPerChannel, adsb.Reserved)
}

//NewAudioStreamBasicDescriptionFromBytes reads AudioStreamBasicDescription from bytes
func NewAudioStreamBasicDescriptionFromBytes(data []byte) (AudioStreamBasicDescription, error) {
	r := bytes.NewReader(data)
	var audioStreamBasicDescription AudioStreamBasicDescription
	err := binary.Read(r, binary.LittleEndian, &audioStreamBasicDescription)
	if err != nil {
		return audioStreamBasicDescription, err
	}
	return audioStreamBasicDescription, nil
}

//SerializeAudioStreamBasicDescription puts an AudioStreamBasicDescription into the given byte array
func (adsb AudioStreamBasicDescription) SerializeAudioStreamBasicDescription(adsbBytes []byte) {
	binary.LittleEndian.PutUint64(adsbBytes, math.Float64bits(adsb.SampleRate))
	var index = 8
	binary.LittleEndian.PutUint32(adsbBytes[index:], AudioFormatIDLpcm)
	index += 4

	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.FormatFlags)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.BytesPerPacket)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.FramesPerPacket)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.BytesPerFrame)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.ChannelsPerFrame)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.BitsPerChannel)
	index += 4
	binary.LittleEndian.PutUint32(adsbBytes[index:], adsb.Reserved)
	index += 4

	binary.LittleEndian.PutUint64(adsbBytes[index:], math.Float64bits(adsb.SampleRate))
	index += 8
	binary.LittleEndian.PutUint64(adsbBytes[index:], math.Float64bits(adsb.SampleRate))

}
