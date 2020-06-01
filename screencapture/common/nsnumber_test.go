package common_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

//I took these from hexdumps
var typeSix = []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x9E, 0x40}
var typeFive = []byte{0x05, 0x2E, 00, 00, 00}
var typeFour = []byte{0x04, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var typeThree = []byte{0x03, 0x1E, 00, 00, 00}

const typeSixDecoded float64 = 1920
const typeFiveDecoded uint32 = 46
const typeFourDecoded uint64 = 5
const typeThreeDecoded uint32 = 30

func TestErrors(t *testing.T) {
	var broken []byte
	broken = make([]byte, len(typeSix))
	copy(broken, typeSix)
	broken[0] = 3
	_, err := common.NewNSNumber(broken)
	assert.Error(t, err)

	broken = make([]byte, len(typeThree))
	copy(broken, typeThree)
	broken[0] = 6
	_, err = common.NewNSNumber(broken)
	assert.Error(t, err)

	broken[0] = 4
	_, err = common.NewNSNumber(broken)
	assert.Error(t, err)

	broken[0] = 56
	_, err = common.NewNSNumber(broken)
	assert.Error(t, err)

	broken = make([]byte, len(typeFive))
	copy(broken, typeFive)
	broken[0] = 134
	_, err = common.NewNSNumber(broken)
	assert.Error(t, err)
}

func TestNumberValue(t *testing.T) {

	float64Num, err := common.NewNSNumber(typeSix)
	if assert.NoError(t, err) {
		assert.Equal(t, typeSixDecoded, float64Num.FloatValue)
	}

	uint32Num, err := common.NewNSNumber(typeFive)
	if assert.NoError(t, err) {
		assert.Equal(t, typeFiveDecoded, uint32Num.IntValue)
	}

	uint64Num, err := common.NewNSNumber(typeFour)
	if assert.NoError(t, err) {
		assert.Equal(t, typeFourDecoded, uint64Num.LongValue)
	}

	uint32Num, err = common.NewNSNumber(typeThree)
	if assert.NoError(t, err) {
		assert.Equal(t, typeThreeDecoded, uint32Num.IntValue)
	}
}

func TestEncoding(t *testing.T) {
	floatNSNumber := common.NewNSNumberFromUFloat64(typeSixDecoded)
	assert.Equal(t, typeSix, floatNSNumber.ToBytes())

	int32NSNumber := common.NewNSNumberFromUInt32(typeThreeDecoded)
	assert.Equal(t, typeThree, int32NSNumber.ToBytes())

	int64NSNumber := common.NewNSNumberFromUInt64(typeFourDecoded)
	assert.Equal(t, typeFour, int64NSNumber.ToBytes())
}
