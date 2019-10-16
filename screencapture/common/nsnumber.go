package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

//NumberValueMagic is vbmn in little endian ascii ==> nmbv
const NumberValueMagic uint32 = 0x6E6D6276

// NSNumber represents a type in the binary protocol used. Type 6 seems to be a float64, type 4 a int64, type 3 a int32.
// I am not sure whether signed or unsigned. They are all in LittleEndian
type NSNumber struct {
	typeSpecifier byte
	//not certain if these are really unsigned
	IntValue   uint32
	LongValue  uint64
	FloatValue float64
}

//NewNSNumberFromUInt32 create NSNumber of type 0x03 with a 4 byte int as value
func NewNSNumberFromUInt32(intValue uint32) NSNumber {
	return NSNumber{typeSpecifier: 03, IntValue: intValue}
}

//NewNSNumberFromUInt64 create NSNumber of type 0x04 with a 8 byte int as value
func NewNSNumberFromUInt64(longValue uint64) NSNumber {
	return NSNumber{typeSpecifier: 04, LongValue: longValue}
}

//NewNSNumberFromUFloat64 create NSNumber of type 0x06 with a 8 byte int as value
func NewNSNumberFromUFloat64(floatValue float64) NSNumber {
	return NSNumber{typeSpecifier: 06, FloatValue: floatValue}
}

//NewNSNumber reads a NSNumber from bytes.
func NewNSNumber(bytes []byte) (NSNumber, error) {
	typeSpecifier := bytes[0]
	switch typeSpecifier {
	case 6:
		if len(bytes) != 9 {
			return NSNumber{}, fmt.Errorf("the NSNumber, type 6 should contain 8 bytes: %s", hex.Dump(bytes))
		}
		value := math.Float64frombits(binary.LittleEndian.Uint64(bytes[1:]))
		return NSNumber{typeSpecifier: typeSpecifier, FloatValue: value}, nil
	case 4:
		if len(bytes) != 9 {
			return NSNumber{}, fmt.Errorf("the NSNumber, type 4 should contain 8 bytes: %s", hex.Dump(bytes))
		}
		value := binary.LittleEndian.Uint64(bytes[1:])
		return NSNumber{typeSpecifier: typeSpecifier, LongValue: value}, nil
	case 3:
		if len(bytes) != 5 {
			return NSNumber{}, fmt.Errorf("the NSNumber, type 3 should contain 4 bytes: %s", hex.Dump(bytes))
		}
		value := binary.LittleEndian.Uint32(bytes[1:])
		return NSNumber{typeSpecifier: typeSpecifier, IntValue: value}, nil
	default:
		return NSNumber{}, fmt.Errorf("unknown NSNumber type %d", typeSpecifier)
	}

}

//ToBytes serializes a NSNumber into a []byte.
//FIXME: remove allocation of array and use one that is passed in instead
func (n NSNumber) ToBytes() []byte {
	switch n.typeSpecifier {
	case 6:
		result := make([]byte, 9)
		binary.LittleEndian.PutUint64(result[1:], math.Float64bits(n.FloatValue))
		result[0] = n.typeSpecifier
		return result
	case 4:
		result := make([]byte, 9)
		binary.LittleEndian.PutUint64(result[1:], n.LongValue)
		result[0] = n.typeSpecifier
		return result
	case 3:
		result := make([]byte, 5)
		binary.LittleEndian.PutUint32(result[1:], n.IntValue)
		result[0] = n.typeSpecifier
		return result
	default:
		//shouldn't happen
		return nil
	}
}

func (n NSNumber) String() string {
	switch n.typeSpecifier {
	case 6:
		return fmt.Sprintf("Float64[%f]", n.FloatValue)
	case 3:
		return fmt.Sprintf("Int32[%d]", n.IntValue)
	case 4:
		return fmt.Sprintf("UInt64[%d]", n.LongValue)
	default:
		return "Invalid Type Specifier"
	}

}
