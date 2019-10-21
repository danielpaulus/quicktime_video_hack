package coremedia_test

import (
	"io/ioutil"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestIntDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/intdict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := coremedia.NewIndexDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(mydict.Entries))
		assert.Equal(t, uint16(49), mydict.Entries[0].Key)
		assert.IsType(t, coremedia.IndexKeyDict{}, mydict.Entries[0].Value)
		nestedDict := mydict.Entries[0].Value.(coremedia.IndexKeyDict)

		assert.Equal(t, 1, len(nestedDict.Entries))
		assert.Equal(t, uint16(105), nestedDict.Entries[0].Key)
		assert.Equal(t, 36, len(nestedDict.Entries[0].Value.([]byte)))

		assert.Equal(t, uint16(52), mydict.Entries[1].Key)
		assert.Equal(t, "H.264", mydict.Entries[1].Value)
	}
}

func TestBooleanEntry(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/bulvalue.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := coremedia.NewStringDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(mydict.Entries))
		assert.Equal(t, "Valeria", mydict.Entries[0].Key)
		assert.Equal(t, true, mydict.Entries[0].Value.(bool))
	}
}

func TestSimpleDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := coremedia.NewStringDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(mydict.Entries))
		assert.Equal(t, "Valeria", mydict.Entries[0].Key)
		assert.Equal(t, true, mydict.Entries[0].Value.(bool))

		assert.Equal(t, "HEVCDecoderSupports444", mydict.Entries[1].Key)
		assert.Equal(t, true, mydict.Entries[1].Value.(bool))

		assert.Equal(t, "DisplaySize", mydict.Entries[2].Key)
		displaySize := mydict.Entries[2].Value.(coremedia.StringKeyDict)
		assert.Equal(t, 2, len(displaySize.Entries))

		assert.Equal(t, "Width", displaySize.Entries[0].Key)
		assert.Equal(t, float64(1920), displaySize.Entries[0].Value.(common.NSNumber).FloatValue)

		assert.Equal(t, "Height", displaySize.Entries[1].Key)
		assert.Equal(t, float64(1200), displaySize.Entries[1].Value.(common.NSNumber).FloatValue)
	}
}

func TestComplexDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/complex_dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := coremedia.NewStringDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(mydict.Entries))
		assert.IsType(t, coremedia.FormatDescriptor{}, mydict.Entries[2].Value)
	}
}

func TestStringFunction(t *testing.T) {
	//TODO: add an assertion
	print(packet.CreateHpa1DeviceInfoDict().String())
	numberDict := coremedia.IndexKeyDict{Entries: make([]coremedia.IndexKeyEntry, 1)}
	numberDict.Entries[0] = coremedia.IndexKeyEntry{
		Key: 5,
		Value: coremedia.FormatDescriptor{
			MediaType:            coremedia.MediaTypeVideo,
			VideoDimensionWidth:  500,
			VideoDimensionHeight: 500,
			Codec:                coremedia.CodecAvc1,
			Extensions:           coremedia.IndexKeyDict{},
		},
	}
	print(numberDict.String())
}
