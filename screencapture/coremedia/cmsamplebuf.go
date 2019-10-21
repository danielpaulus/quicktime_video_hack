package coremedia

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
)

//CMItemCount is a simple typedef to int to be a bit closer to MacOS/iOS
type CMItemCount = int

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

//CMSampleTimingInfo is a simple struct containing 3 CMtimes: Duration, PresentationTimeStamp and DecodeTimeStamp
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

func (info CMSampleTimingInfo) String() string {
	return fmt.Sprintf("{Duration:%s, PresentationTS:%s, DecodeTS:%s}",
		info.Duration, info.PresentationTimeStamp, info.DecodeTimeStamp)
}

//CMSampleBuffer represents the CoreMedia class used to exchange AV SampleData and contains meta information like timestamps or
//optional FormatDescriptors
type CMSampleBuffer struct {
	OutputPresentationTimestamp CMTime
	FormatDescription           FormatDescriptor
	HasFormatDescription        bool
	NumSamples                  CMItemCount          //nsmp
	SampleTimingInfoArray       []CMSampleTimingInfo //stia
	SampleData                  []byte
	SampleSizes                 []int
	Attachments                 IndexKeyDict //satt
	Sary                        IndexKeyDict //sary
	MediaType                   uint32
}

func (buffer CMSampleBuffer) String() string {
	var fdscString string
	if buffer.HasFormatDescription {
		fdscString = buffer.FormatDescription.String()
	} else {
		fdscString = "none"
	}
	if buffer.MediaType == MediaTypeVideo {
		return fmt.Sprintf("{OutputPresentationTS:%s, NumSamples:%d, Nalus:%s, fdsc:%s, attach:%s, sary:%s, SampleTimingInfoArray:%s}",
			buffer.OutputPresentationTimestamp.String(), buffer.NumSamples, GetNaluDetails(buffer.SampleData),
			fdscString, buffer.Attachments.String(), buffer.Sary.String(), buffer.SampleTimingInfoArray[0].String())
	}
	return fmt.Sprintf("{OutputPresentationTS:%s, NumSamples:%d, SampleSize:%d, fdsc:%s}",
		buffer.OutputPresentationTimestamp.String(), buffer.NumSamples, buffer.SampleSizes[0],
		fdscString)
}

//NewCMSampleBufferFromBytesAudio parses a CMSampleBuffer containing audio data.
func NewCMSampleBufferFromBytesAudio(data []byte) (CMSampleBuffer, error) {
	return NewCMSampleBufferFromBytes(data, MediaTypeSound)
}

//NewCMSampleBufferFromBytesVideo parses a CMSampleBuffer containing audio video.
func NewCMSampleBufferFromBytesVideo(data []byte) (CMSampleBuffer, error) {
	return NewCMSampleBufferFromBytes(data, MediaTypeVideo)
}

//NewCMSampleBufferFromBytes parses a CMSampleBuffer from a []byte assuming it begins with a 4 byte length and the 4byte magic int "sbuf"
func NewCMSampleBufferFromBytes(data []byte, mediaType uint32) (CMSampleBuffer, error) {
	var sbuffer CMSampleBuffer
	sbuffer.MediaType = mediaType
	sbuffer.HasFormatDescription = false
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

	length, remainingBytes, err = common.ParseLengthAndMagic(remainingBytes, sdat)
	if err != nil {
		return sbuffer, err
	}
	sbuffer.SampleData = remainingBytes[:length-8]
	length, remainingBytes, err = common.ParseLengthAndMagic(remainingBytes[length-8:], nsmp)
	if err != nil {
		return sbuffer, err
	}
	if length != 12 {
		return sbuffer, fmt.Errorf("invalid length for nsmp %d, should be 12", length)
	}
	sbuffer.NumSamples = int(binary.LittleEndian.Uint32(remainingBytes))

	sbuffer.SampleSizes, remainingBytes, err = parseSampleSizeArray(remainingBytes[4:])
	if err != nil {
		return sbuffer, err
	}

	//audio buffers usually end after samplesize
	if len(remainingBytes) == 0 {
		return sbuffer, nil
	}
	if binary.LittleEndian.Uint32(remainingBytes[4:]) == FormatDescriptorMagic {
		sbuffer.HasFormatDescription = true
		fdscLength := binary.LittleEndian.Uint32(remainingBytes)
		sbuffer.FormatDescription, err = NewFormatDescriptorFromBytes(remainingBytes[:fdscLength])
		if err != nil {
			return sbuffer, err
		}
		remainingBytes = remainingBytes[fdscLength:]
	}
	//audio buffers usually end after samplesize
	if len(remainingBytes) == 0 {
		return sbuffer, nil
	}

	attachmentsLength := binary.LittleEndian.Uint32(remainingBytes)
	sbuffer.Attachments, err = NewIndexDictFromBytesWithCustomMarker(remainingBytes[:attachmentsLength], satt)
	if err != nil {
		return sbuffer, err
	}
	remainingBytes = remainingBytes[attachmentsLength:]
	saryLength := binary.LittleEndian.Uint32(remainingBytes)
	if binary.LittleEndian.Uint32(remainingBytes[4:]) != sary {
		return sbuffer, fmt.Errorf("wrong magic, expected sary got:%x", remainingBytes[4:8])
	}
	sbuffer.Sary, err = NewIndexDictFromBytes(remainingBytes[8:saryLength])
	if err != nil {
		return sbuffer, err
	}
	if len(remainingBytes[saryLength:]) != 0 {
		return sbuffer, fmt.Errorf("CmSampleBuf should have been read completely but still contains bytes: %x", remainingBytes[saryLength:])
	}
	return sbuffer, nil
}

func parseSampleSizeArray(data []byte) ([]int, []byte, error) {
	ssizLength, _, err := common.ParseLengthAndMagic(data, ssiz)
	if err != nil {
		return nil, nil, err
	}
	ssizLength -= 8
	numEntries, modulus := ssizLength/4, ssizLength%4
	if modulus != 0 {
		return nil, nil, fmt.Errorf("error parsing samplesizearray, too many bytes: %d", modulus)
	}
	result := make([]int, numEntries)
	data = data[8:]
	for i := 0; i < numEntries; i++ {
		index := 4 * i
		result[i] = int(binary.LittleEndian.Uint32(data[index+i*4:]))
	}
	return result, data[ssizLength:], nil
}

func parseStia(data []byte) ([]CMSampleTimingInfo, []byte, error) {
	stiaLength, _, err := common.ParseLengthAndMagic(data, stia)
	if err != nil {
		return nil, nil, err
	}
	stiaLength -= 8

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
