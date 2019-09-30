package dict

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

const Magic int = 0x76626D6E

// Type 6 seems to be a float64, type 4 a int64, type 3 a int32.
// I am not sure whether signed or unsigned. They are all in LittleEndian
type NSNumber struct {
	typeSpecifier byte
	//not certain if these are really unsigned
	IntValue   uint32
	LongValue  uint64
	FloatValue float64
}

//Read what I assume is a NSNumber from bytes
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
