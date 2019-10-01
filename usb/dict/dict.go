package dict

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	KeyValuePairMagic uint32 = 0x6B657976 //keyv - vyek
	StringKey         uint32 = 0x7374726B //strk - krts
	IntKey            uint32 = 0x6964786B //idxk - kxdi
	BooleanValueMagic uint32 = 0x62756C76 //bulv - vlub
	DictionaryMagic   uint32 = 0x64696374 //dict - tcid
	DataValueMagic    uint32 = 0x64617476 //datv - vtad
	StringValueMagic  uint32 = 0x73747276 //strv - vrts
)

type StringKeyDict struct {
	Entries []StringKeyEntry
}
type StringKeyEntry struct {
	Key   string
	Value interface{}
}

type IntKeyDict struct {
	Entries []IntKeyEntry
}
type IntKeyEntry struct {
	Key   uint16
	Value interface{}
}

func NewIntDictFromBytes(data []byte) (IntKeyDict, error) {
	return NewIntDictFromBytesWithCustomMarker(data, DictionaryMagic)
}

func NewIntDictFromBytesWithCustomMarker(data []byte, magic uint32) (IntKeyDict, error) {
	_, remainingBytes, err := parseLengthAndMagic(data, magic)
	if err != nil {
		return IntKeyDict{}, err
	}
	var slice = remainingBytes
	dict := IntKeyDict{}
	for len(slice) != 0 {
		keyValuePairLength := binary.LittleEndian.Uint32(slice)
		if int(keyValuePairLength) > len(slice) {
			return dict, fmt.Errorf("invalid dict: %s", hex.Dump(data))
		}
		keyValuePair := slice[8:keyValuePairLength]
		intDictEntry, err := parseIntDictEntry(keyValuePair)
		if err != nil {
			return dict, err
		}
		dict.Entries = append(dict.Entries, intDictEntry)
		slice = slice[keyValuePairLength:]
	}
	return dict, nil
}

func NewStringDictFromBytes(data []byte) (StringKeyDict, error) {
	_, remainingBytes, err := parseLengthAndMagic(data, DictionaryMagic)
	if err != nil {
		return StringKeyDict{}, err
	}

	var slice = remainingBytes
	dict := StringKeyDict{}
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

func parseIntDictEntry(bytes []byte) (IntKeyEntry, error) {
	key, remainingBytes, err := parseIntKey(bytes)
	if err != nil {
		return IntKeyEntry{}, err
	}
	value, err := parseValue(remainingBytes)
	if err != nil {
		return IntKeyEntry{}, err
	}
	return IntKeyEntry{Key: key, Value: value}, nil
}

func parseEntry(bytes []byte) (StringKeyEntry, error) {
	key, remainingBytes, err := parseKey(bytes)
	if err != nil {
		return StringKeyEntry{}, err
	}
	value, err := parseValue(remainingBytes)
	if err != nil {
		return StringKeyEntry{}, err
	}
	return StringKeyEntry{Key: key, Value: value}, nil
}

func parseKey(bytes []byte) (string, []byte, error) {
	keyLength, _, err := parseLengthAndMagic(bytes, StringKey)
	if err != nil {
		return "", nil, err
	}
	key := string(bytes[8:keyLength])
	return key, bytes[keyLength:], nil
}

func parseIntKey(bytes []byte) (uint16, []byte, error) {
	keyLength := binary.LittleEndian.Uint32(bytes)
	if len(bytes) < int(keyLength) {
		return 0, nil, fmt.Errorf("invalid key data length, cannot parse string %s", hex.Dump(bytes))
	}
	magic := binary.LittleEndian.Uint32(bytes[4:])
	if IntKey != magic {
		return 0, nil, fmt.Errorf("invalid key magic:%x, cannot parse string %s", magic, hex.Dump(bytes))
	}
	key := binary.LittleEndian.Uint16(bytes[8:])
	return key, bytes[keyLength:], nil
}

func parseValue(bytes []byte) (interface{}, error) {
	valueLength := binary.LittleEndian.Uint32(bytes)
	if len(bytes) < int(valueLength) {
		return nil, fmt.Errorf("invalid value data length, cannot parse %s", hex.Dump(bytes))
	}
	magic := binary.LittleEndian.Uint32(bytes[4:])
	switch magic {
	case StringValueMagic:
		return string(bytes[8:valueLength]), nil
	case DataValueMagic:
		return bytes[8:valueLength], nil
	case BooleanValueMagic:
		return bytes[8] == 1, nil
	case NumberValueMagic:
		return NewNSNumber(bytes[8:])
	case DictionaryMagic:
		//FIXME: that is a lazy implementation, improve please
		dict, err := NewStringDictFromBytes(bytes)
		if err != nil {
			return NewIntDictFromBytes(bytes)
		}
		return dict, nil
	case FormatDescriptorMagic:
		return NewFormatDescriptorFromBytes(bytes)
	default:
		unknownMagic := string(bytes[4:8])
		return nil, fmt.Errorf("unknown dictionary magic type:%s (0x%x), cannot parse value %s", unknownMagic, magic, hex.Dump(bytes))
	}
}

func (ikd IntKeyDict) getValue(index uint16) interface{} {
	for _, entry := range ikd.Entries {
		if entry.Key == index {
			return entry.Value
		}
	}
	return nil
}
