package coremedia

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/usb/common"
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

type CMItemCount = uint32

//https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.9.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMSampleBuffer.h
const (
	sbuf uint32 = 0x73627566 //the cmsamplebuf and only content of feed asyns
	opts uint32 = 0x6F707473 //output presentation timestamp?
	stia uint32 = 0x73746961 //sampleTimingInfoArray
	sdat uint32 = 0x73646174 //the nalu
	satt uint32 = 0x73617474 //indexkey dict with only number values, CMSampleBufferGetSampleAttachmentsArray
	sary uint32 = 0x73617279 //some dict with index and one boolean
	ssiz uint32 = 0x7373697A //samplesize in bytes, size of what is contained in sdat, sample size array i think
	nsmp uint32 = 0x6E736D70 //numsample so you know how many things are in the arrays

	cmSampleTimingInfoLength = 3 * CMTimeLengthInBytes
)

type CMSampleTimingInfo struct {
	Duration CMTime /*! @field duration
	The duration of the sample. If a single struct applies to
	each of the samples, they all will have this duration. */
	PresentationTimeStamp CMTime /*! @field presentationTimeStamp
	The time at which the sample will be presented. If a single
	struct applies to each of the samples, this is the presentationTime of the
	first sample. The presentationTime of subsequent samples will be derived by
	repeatedly adding the sample duration. */
	DecodeTimeStamp CMTime /*! @field decodeTimeStamp
	The time at which the sample will be decoded. If the samples
	are in presentation order, this must be set to kCMTimeInvalid. */
}

type CMSampleBuffer struct {
	OutputPresentationTimestamp CMTime
	FormatDescription           dict.FormatDescriptor
	NumSamples                  CMItemCount          //nsmp
	SampleTimingInfoArray       []CMSampleTimingInfo //stia
	SampleData                  []byte
}

func NewCMSampleBufferFromBytes(data []byte) (CMSampleBuffer, error) {
	var sbuffer CMSampleBuffer
	length, remainingBytes, err := common.ParseLengthAndMagic(data, sbuf)
	if err != nil {
		return sbuffer, err
	}
	if length > len(data) {
		return sbuffer, fmt.Errorf("less data (%d bytes) in buffer than expected (%d bytes)", len(data), length)
	}

	_, remainingBytes, err = common.ParseLengthAndMagic(remainingBytes, opts)
	if err != nil {
		return sbuffer, err
	}
	cmtime, err := NewCMTimeFromBytes(remainingBytes)
	if err != nil {
		return sbuffer, err
	}
	sbuffer.OutputPresentationTimestamp = cmtime
	sbuffer.SampleTimingInfoArray, remainingBytes, err = parseStia(remainingBytes[24:])
	if err != nil {
		return sbuffer, err
	}
	return sbuffer, nil
}

func parseStia(data []byte) ([]CMSampleTimingInfo, []byte, error) {
	stiaLength := int(binary.LittleEndian.Uint32(data) - 8)

	numEntries, modulus := stiaLength/cmSampleTimingInfoLength, stiaLength%cmSampleTimingInfoLength
	if modulus != 0 {
		return nil, nil, fmt.Errorf("error parsing stia, too many bytes: %d", modulus)
	}
	result := make([]CMSampleTimingInfo, numEntries)
	data = data[8:]
	for i := 0; i < numEntries; i++ {
		index := i * cmSampleTimingInfoLength
		duration, err := NewCMTimeFromBytes(data[index:])
		if err != nil {
			return nil, nil, err
		}
		presentationTimeStamp, err := NewCMTimeFromBytes(data[CMTimeLengthInBytes+index:])
		if err != nil {
			return nil, nil, err
		}
		decodeTimeStamp, err := NewCMTimeFromBytes(data[2*CMTimeLengthInBytes+index:])
		if err != nil {
			return nil, nil, err
		}

		result[i] = CMSampleTimingInfo{
			Duration:              duration,
			PresentationTimeStamp: presentationTimeStamp,
			DecodeTimeStamp:       decodeTimeStamp,
		}
	}
	return result, data[stiaLength:], nil
}
