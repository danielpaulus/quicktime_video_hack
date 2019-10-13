package dict_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
	"github.com/danielpaulus/quicktime_video_hack/usb/messages"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestIntDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/intdict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := dict.NewIndexDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(mydict.Entries))
		assert.Equal(t, uint16(49), mydict.Entries[0].Key)
		assert.IsType(t, dict.IndexKeyDict{}, mydict.Entries[0].Value)
		nestedDict := mydict.Entries[0].Value.(dict.IndexKeyDict)

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
	mydict, err := dict.NewStringDictFromBytes(dat)
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
	mydict, err := dict.NewStringDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(mydict.Entries))
		assert.Equal(t, "Valeria", mydict.Entries[0].Key)
		assert.Equal(t, true, mydict.Entries[0].Value.(bool))

		assert.Equal(t, "HEVCDecoderSupports444", mydict.Entries[1].Key)
		assert.Equal(t, true, mydict.Entries[1].Value.(bool))

		assert.Equal(t, "DisplaySize", mydict.Entries[2].Key)
		displaySize := mydict.Entries[2].Value.(dict.StringKeyDict)
		assert.Equal(t, 2, len(displaySize.Entries))

		assert.Equal(t, "Width", displaySize.Entries[0].Key)
		assert.Equal(t, float64(1920), displaySize.Entries[0].Value.(dict.NSNumber).FloatValue)

		assert.Equal(t, "Height", displaySize.Entries[1].Key)
		assert.Equal(t, float64(1200), displaySize.Entries[1].Value.(dict.NSNumber).FloatValue)
	}
}

func TestComplexDict(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/complex_dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	mydict, err := dict.NewStringDictFromBytes(dat)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(mydict.Entries))
		assert.IsType(t, dict.FormatDescriptor{}, mydict.Entries[2].Value)
	}
}

func TestStringFunction(t *testing.T) {
	//TODO: add an assertion
	print(messages.CreateHpa1DeviceInfoDict().String())
	numberDict := dict.IndexKeyDict{Entries: make([]dict.IndexKeyEntry, 1)}
	numberDict.Entries[0] = dict.IndexKeyEntry{
		Key: 5,
		Value: dict.FormatDescriptor{
			MediaType:            dict.MediaTypeVideo,
			VideoDimensionWidth:  500,
			VideoDimensionHeight: 500,
			Codec:                dict.CodecAvc1,
			Extensions:           dict.IndexKeyDict{},
		},
	}
	print(numberDict.String())
}
