package dict

import (
	"encoding/binary"
	"fmt"
)

//Those are the markers found in the hex dumps.
//For convenience I have added the ASCII representation as a comment
//in normal byte order and reverse byteorder (so you can find them in the hex dumps)
// Note: I have just guessed what the names could be from the marker ascii, I could be wrong ;-)
const (
	FormatDescriptorMagic uint32 = 0x66647363 //fdsc - csdf
	MediaTypeVideo        uint32 = 0x76696465 //vide - ediv
	MediaTypeMagic        uint32 = 0x6D646961 //mdia - aidm
	VideoDimensionMagic   uint32 = 0x7664696D //vdim - midv
	CodecMagic            uint32 = 0x636F6463 //codc - cdoc
	CodecAvc1             uint32 = 0x61766331 //avc1 - 1cva
	ExtensionMagic        uint32 = 0x6578746E //extn - ntxe
)

//Seems like a CMFormatDescription
//https://developer.apple.com/documentation/coremedia/cmformatdescription
//https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.9.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMFormatDescription.h
type FormatDescriptor struct {
	MediaType            uint32
	VideoDimensionWidth  uint32
	VideoDimensionHeight uint32
	Codec                uint32
	Extensions           IndexKeyDict
}

func NewFormatDescriptorFromBytes(data []byte) (FormatDescriptor, error) {

	_, remainingBytes, err := parseLengthAndMagic(data, FormatDescriptorMagic)
	if err != nil {
		return FormatDescriptor{}, err
	}
	mediaType, remainingBytes, err := parseMediaType(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	videoDimensionWidth, videoDimensionHeight, remainingBytes, err := parseVideoDimension(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	codec, remainingBytes, err := parseCodec(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	extensions, err := NewIndexDictFromBytesWithCustomMarker(remainingBytes, ExtensionMagic)

	return FormatDescriptor{
		MediaType:            mediaType,
		VideoDimensionHeight: videoDimensionHeight,
		VideoDimensionWidth:  videoDimensionWidth,
		Codec:                codec,
		Extensions:           extensions, //doc on extensions at the bottom of: https://developer.apple.com/documentation/coremedia/cmformatdescription?language=objc
	}, nil
}

func parseCodec(bytes []byte) (uint32, []byte, error) {
	length, _, err := parseLengthAndMagic(bytes, CodecMagic)
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
	length, _, err := parseLengthAndMagic(bytes, VideoDimensionMagic)
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
	length, _, err := parseLengthAndMagic(bytes, MediaTypeMagic)
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
	return fmt.Sprintf(
		"FormatDescriptor:\n\t MediaType %s \n\t VideoDimension:(%dx%d) \n\t Codec:%s \n\t Extensions:%s \n",
		readableMediaType(fdsc.MediaType), fdsc.VideoDimensionWidth, fdsc.VideoDimensionHeight,
		readableCodec(fdsc.Codec), fdsc.Extensions.String())
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
	return fmt.Sprintf("Unknown(%x)", mediaType)
}
