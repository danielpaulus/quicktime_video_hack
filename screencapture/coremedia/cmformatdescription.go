package coremedia

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/sirupsen/logrus"
)

//Those are the markers found in the hex dumps.
//For convenience I have added the ASCII representation as a comment
//in normal byte order and reverse byteorder (so you can find them in the hex dumps)
// Note: I have just guessed what the names could be from the marker ascii, I could be wrong ;-)
const (
	FormatDescriptorMagic            uint32 = 0x66647363 //fdsc - csdf
	MediaTypeVideo                   uint32 = 0x76696465 //vide - ediv
	MediaTypeSound                   uint32 = 0x736F756E //nuos - soun
	MediaTypeMagic                   uint32 = 0x6D646961 //mdia - aidm
	VideoDimensionMagic              uint32 = 0x7664696D //vdim - midv
	CodecMagic                       uint32 = 0x636F6463 //codc - cdoc
	CodecAvc1                        uint32 = 0x61766331 //avc1 - 1cva
	ExtensionMagic                   uint32 = 0x6578746E //extn - ntxe
	AudioStreamBasicDescriptionMagic uint32 = 0x61736264 //asdb - dbsa
)

//FormatDescriptor is actually a CMFormatDescription
//https://developer.apple.com/documentation/coremedia/cmformatdescription
//https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.9.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMFormatDescription.h
type FormatDescriptor struct {
	MediaType            uint32
	VideoDimensionWidth  uint32
	VideoDimensionHeight uint32
	Codec                uint32
	Extensions           IndexKeyDict
	//PPS contains bytes of the Picture Parameter Set h264 NALu
	PPS []byte
	//SPS contains bytes of the Picture Parameter Set h264 NALu
	SPS                         []byte
	AudioStreamBasicDescription AudioStreamBasicDescription
}

//NewFormatDescriptorFromBytes parses a CMFormatDescription from bytes
func NewFormatDescriptorFromBytes(data []byte) (FormatDescriptor, error) {

	_, remainingBytes, err := common.ParseLengthAndMagic(data, FormatDescriptorMagic)
	if err != nil {
		return FormatDescriptor{}, err
	}
	mediaType, remainingBytes, err := parseMediaType(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	if mediaType == MediaTypeSound {
		return parseSoundFdsc(remainingBytes)
	}
	return parseVideoFdsc(remainingBytes)
}
func parseSoundFdsc(remainingBytes []byte) (FormatDescriptor, error) {

	length, _, err := common.ParseLengthAndMagic(remainingBytes, AudioStreamBasicDescriptionMagic)
	if err != nil {
		return FormatDescriptor{}, err
	}

	asdb, err := NewAudioStreamBasicDescriptionFromBytes(remainingBytes[8:length])
	if err != nil {
		return FormatDescriptor{}, err
	}

	return FormatDescriptor{
		MediaType:                   MediaTypeSound,
		AudioStreamBasicDescription: asdb,
	}, nil
}
func parseVideoFdsc(remainingBytes []byte) (FormatDescriptor, error) {
	videoDimensionWidth, videoDimensionHeight, remainingBytes, err := parseVideoDimension(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	codec, remainingBytes, err := parseCodec(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	extensions, err := NewIndexDictFromBytesWithCustomMarker(remainingBytes, ExtensionMagic)
	if err != nil {
		return FormatDescriptor{}, err
	}

	pps, sps := extractPPS(extensions)
	return FormatDescriptor{
		MediaType:            MediaTypeVideo,
		VideoDimensionHeight: videoDimensionHeight,
		VideoDimensionWidth:  videoDimensionWidth,
		Codec:                codec,
		Extensions:           extensions, //doc on extensions at the bottom of: https://developer.apple.com/documentation/coremedia/cmformatdescription?language=objc
		PPS:                  pps,
		SPS:                  sps,
	}, nil
}

func extractPPS(dict IndexKeyDict) ([]byte, []byte) {
	val, err := dict.getValue(49)
	if err != nil {
		logrus.Error("FDSC did not contain PPS/SPS")
		return make([]byte, 0), make([]byte, 0)
	}
	val, err = val.(IndexKeyDict).getValue(105)
	if err != nil {
		logrus.Error("FDSC did not contain PPS/SPS")
		return make([]byte, 0), make([]byte, 0)
	}
	data := val.([]byte)
	ppsLength := data[7]
	pps := data[8 : 8+ppsLength]
	spsLength := data[10+ppsLength]
	sps := data[11+ppsLength : 11+ppsLength+spsLength]
	return pps, sps
}

func parseCodec(bytes []byte) (uint32, []byte, error) {
	length, _, err := common.ParseLengthAndMagic(bytes, CodecMagic)
	if err != nil {
		return 0, nil, err
	}
	if length != 12 {
		return 0, nil, fmt.Errorf("invalid length for codec: %d", length)
	}
	codec := binary.LittleEndian.Uint32(bytes[8:])
	return codec, bytes[length:], nil
}

func parseVideoDimension(bytes []byte) (uint32, uint32, []byte, error) {
	length, _, err := common.ParseLengthAndMagic(bytes, VideoDimensionMagic)
	if err != nil {
		return 0, 0, nil, err
	}
	if length != 16 {
		return 0, 0, nil, fmt.Errorf("invalid length for video dimension: %d", length)
	}
	width := binary.LittleEndian.Uint32(bytes[8:])
	height := binary.LittleEndian.Uint32(bytes[12:])
	return width, height, bytes[length:], nil
}

func parseMediaType(bytes []byte) (uint32, []byte, error) {
	length, _, err := common.ParseLengthAndMagic(bytes, MediaTypeMagic)
	if err != nil {
		return 0, nil, err
	}
	if length != 12 {
		return 0, nil, fmt.Errorf("invalid length for media type: %d", length)
	}
	mediaType := binary.LittleEndian.Uint32(bytes[8:])
	return mediaType, bytes[length:], nil
}

func (fdsc FormatDescriptor) String() string {
	if fdsc.MediaType == MediaTypeVideo {
		return fmt.Sprintf(
			"fdsc:{MediaType:%s, VideoDimension:(%dx%d), Codec:%s, PPS:%x, SPS:%x, Extensions:%s}",
			readableMediaType(fdsc.MediaType), fdsc.VideoDimensionWidth, fdsc.VideoDimensionHeight,
			readableCodec(fdsc.Codec), fdsc.PPS, fdsc.SPS, fdsc.Extensions.String())
	}
	return fmt.Sprintf(
		"fdsc:{MediaType:%s, AudioStreamBasicDescription: %s}", readableMediaType(fdsc.MediaType), fdsc.AudioStreamBasicDescription.String())
}

func readableCodec(codec uint32) string {
	if codec == CodecAvc1 {
		return "AVC-1"
	}
	return fmt.Sprintf("Unknown(%x)", codec)
}

func readableMediaType(mediaType uint32) string {
	if mediaType == MediaTypeVideo {
		return "Video"
	}
	if mediaType == MediaTypeSound {
		return "Sound"
	}
	return fmt.Sprintf("Unknown(%x)", mediaType)
}
