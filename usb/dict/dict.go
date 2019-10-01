package dict

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const strk = 0x7374726B

//bulv in ascii, vlub in little endian
const BooleanValueMagic = 0x62756C76

//dict or tcid
const DictMagicValue = 0x64696374

type Dict struct {
	Entries []DictEntry
}
type DictEntry struct {
	Key   string
	Value interface{}
}

func NewDictFromBytes(data []byte) (Dict, error) {
	dictLength := binary.LittleEndian.Uint32(data)
	if int(dictLength) > len(data) {
		return Dict{}, fmt.Errorf("invalid dict: %s", hex.Dump(data))
	}
	magic := binary.LittleEndian.Uint32(data[4:])
	if DictMagicValue != magic {
		unknownMagic := string(data[4:8])
		return Dict{}, fmt.Errorf("invalid dict magic:%s (0x%x), cannot parse dict %s", unknownMagic, magic, hex.Dump(data))
	}

	var slice = data[8:]
	dict := Dict{}
	for len(slice) != 0 {
		keyValuePairLength := binary.LittleEndian.Uint32(slice)
		if int(keyValuePairLength) > len(slice) {
			return dict, fmt.Errorf("invalid dict: %s", hex.Dump(data))
		}
		keyValuePair := slice[8:keyValuePairLength]
		parseDictEntry, err := parseEntry(keyValuePair)
		if err != nil {
			return dict, err
		}
		dict.Entries = append(dict.Entries, parseDictEntry)
		slice = slice[keyValuePairLength:]
	}
	return dict, nil
}

func parseEntry(bytes []byte) (DictEntry, error) {
	key, remainingBytes, err := parseKey(bytes)
	if err != nil {
		return DictEntry{}, err
	}
	value, err := parseValue(remainingBytes)
	if err != nil {
		return DictEntry{}, err
	}
	return DictEntry{Key: key, Value: value}, nil
}

func parseKey(bytes []byte) (string, []byte, error) {
	keyLength := binary.LittleEndian.Uint32(bytes)
	if len(bytes) < int(keyLength) {
		return "", nil, fmt.Errorf("invalid key data length, cannot parse string %s", hex.Dump(bytes))
	}
	magic := binary.LittleEndian.Uint32(bytes[4:])
	if strk != magic {
		return "", nil, fmt.Errorf("invalid key magic:%x, cannot parse string %s", magic, hex.Dump(bytes))
	}
	key := string(bytes[8:keyLength])
	return key, bytes[keyLength:], nil
}

func parseValue(bytes []byte) (interface{}, error) {
	valueLength := binary.LittleEndian.Uint32(bytes)
	if len(bytes) < int(valueLength) {
		return nil, fmt.Errorf("invalid value data length, cannot parse %s", hex.Dump(bytes))
	}
	magic := binary.LittleEndian.Uint32(bytes[4:])
	switch int(magic) {
	case BooleanValueMagic:
		return bytes[8] == 1, nil
	case NumberValueMagic:
		return NewNSNumber(bytes[8:])
	case DictMagicValue:
		return NewDictFromBytes(bytes)
	case FormatDescriptorMagic:
		return NewFormatDescriptorFromBytes(bytes)
	default:
		unknownMagic := string(bytes[4:8])
		return nil, fmt.Errorf("unknown dictionary magic type:%s (0x%x), cannot parse value %s", unknownMagic, magic, hex.Dump(bytes))
	}
	return bytes, nil
}
