package coremedia

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
)

// Dictionary related magic marker constants
const (
	KeyValuePairMagic uint32 = 0x6B657976 //keyv - vyek
	StringKey         uint32 = 0x7374726B //strk - krts
	IntKey            uint32 = 0x6964786B //idxk - kxdi
	BooleanValueMagic uint32 = 0x62756C76 //bulv - vlub
	DictionaryMagic   uint32 = 0x64696374 //dict - tcid
	DataValueMagic    uint32 = 0x64617476 //datv - vtad
	StringValueMagic  uint32 = 0x73747276 //strv - vrts
)

//StringKeyDict a dictionary that uses strings as keys with an array of StringKeyEntry
type StringKeyDict struct {
	Entries []StringKeyEntry
}

//StringKeyEntry a pair of a string key and an arbitrary value
type StringKeyEntry struct {
	Key   string
	Value interface{}
}

//IndexKeyDict a dictionary that uses uint16 as keys with an array of IndexKeyEntry
type IndexKeyDict struct {
	Entries []IndexKeyEntry
}

//IndexKeyEntry is a pair of a uint16 key and an arbitrary value.
type IndexKeyEntry struct {
	Key   uint16
	Value interface{}
}

//NewIndexDictFromBytes creates a new dictionary assuming the byte array starts with the 4 byte length of the dictionary followed by "dict" as the magic marker
func NewIndexDictFromBytes(data []byte) (IndexKeyDict, error) {
	return NewIndexDictFromBytesWithCustomMarker(data, DictionaryMagic)
}

//NewIndexDictFromBytesWithCustomMarker creates a new dictionary assuming the byte array starts with the 4 byte length of the dictionary followed by magic as the magic marker
func NewIndexDictFromBytesWithCustomMarker(data []byte, magic uint32) (IndexKeyDict, error) {
	_, remainingBytes, err := common.ParseLengthAndMagic(data, magic)
	if err != nil {
		return IndexKeyDict{}, err
	}
	var slice = remainingBytes
	dict := IndexKeyDict{}
	for len(slice) != 0 {
		keyValuePairLength, _, err := common.ParseLengthAndMagic(slice, KeyValuePairMagic)
		if err != nil {
			return IndexKeyDict{}, err
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

//NewStringDictFromBytes creates a new dictionary assuming the byte array starts with the 4 byte length of the dictionary followed by "dict" as the magic marker
func NewStringDictFromBytes(data []byte) (StringKeyDict, error) {
	_, remainingBytes, err := common.ParseLengthAndMagic(data, DictionaryMagic)
	if err != nil {
		return StringKeyDict{}, err
	}

	var slice = remainingBytes
	dict := StringKeyDict{}
	for len(slice) != 0 {
		keyValuePairLength, _, err := common.ParseLengthAndMagic(slice, KeyValuePairMagic)
		if err != nil {
			return StringKeyDict{}, err
		}
		keyValuePairData := slice[8:keyValuePairLength]
		parseDictEntry, err := parseEntry(keyValuePairData)
		if err != nil {
			return dict, err
		}
		dict.Entries = append(dict.Entries, parseDictEntry)
		slice = slice[keyValuePairLength:]
	}
	return dict, nil
}

func parseIntDictEntry(bytes []byte) (IndexKeyEntry, error) {
	key, remainingBytes, err := parseIntKey(bytes)
	if err != nil {
		return IndexKeyEntry{}, err
	}
	value, err := parseValue(remainingBytes)
	if err != nil {
		return IndexKeyEntry{}, err
	}
	return IndexKeyEntry{Key: key, Value: value}, nil
}

//ParseKeyValueEntry parses a byte array into a StringKeyEntry assuming the array starts with a 4 byte length followed by the "keyv" magic
func ParseKeyValueEntry(data []byte) (StringKeyEntry, error) {
	keyValuePairLength, _, err := common.ParseLengthAndMagic(data, KeyValuePairMagic)
	if err != nil {
		return StringKeyEntry{}, err
	}
	keyValuePairData := data[8:keyValuePairLength]
	return parseEntry(keyValuePairData)
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
	keyLength, _, err := common.ParseLengthAndMagic(bytes, StringKey)
	if err != nil {
		return "", nil, err
	}
	key := string(bytes[8:keyLength])
	return key, bytes[keyLength:], nil
}

func parseIntKey(bytes []byte) (uint16, []byte, error) {
	keyLength, _, err := common.ParseLengthAndMagic(bytes, IntKey)
	if err != nil {
		return 0, nil, err
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
	case common.NumberValueMagic:
		return common.NewNSNumber(bytes[8:])
	case DictionaryMagic:
		//FIXME: that is a lazy implementation, improve please
		dict, err := NewStringDictFromBytes(bytes)
		if err != nil {
			return NewIndexDictFromBytes(bytes)
		}
		return dict, nil
	case FormatDescriptorMagic:
		return NewFormatDescriptorFromBytes(bytes)
	default:
		unknownMagic := string(bytes[4:8])
		return nil, fmt.Errorf("unknown dictionary magic type:%s (0x%x), cannot parse value %s", unknownMagic, magic, hex.Dump(bytes))
	}
}

func (dt StringKeyDict) String() string {
	sb := strings.Builder{}
	for _, e := range dt.Entries {
		appendEntry(&sb, e)
	}
	return fmt.Sprintf("StringKeyDict:[%s]", sb.String())
}

func (dt IndexKeyDict) String() string {
	sb := strings.Builder{}
	for _, e := range dt.Entries {
		appendIndexEntry(&sb, e)
	}
	return fmt.Sprintf("IndexKeyDict:[%s]", sb.String())
}

func appendIndexEntry(builder *strings.Builder, entry IndexKeyEntry) {
	builder.WriteString("{")
	builder.WriteString(fmt.Sprintf("%d", entry.Key))
	builder.WriteString(" : ")
	valueToString(builder, entry.Value)
	builder.WriteString("},")
}

func appendEntry(builder *strings.Builder, entry StringKeyEntry) {
	builder.WriteString("{")
	builder.WriteString(entry.Key)
	builder.WriteString(" : ")
	valueToString(builder, entry.Value)
	builder.WriteString("},")
}

func valueToString(builder *strings.Builder, value interface{}) {
	switch value := value.(type) {
	case common.NSNumber:
		builder.WriteString(value.String())
	case StringKeyDict:
		builder.WriteString(value.String())
	case []byte:
		builder.WriteString(fmt.Sprintf("0x%x", value))
	case FormatDescriptor:
		builder.WriteString(value.String())
	default:
		builder.WriteString(fmt.Sprintf("%s", value))
	}
}

func (dt IndexKeyDict) getValue(index uint16) (interface{}, error) {
	for _, entry := range dt.Entries {
		if entry.Key == index {
			return entry.Value, nil
		}
	}
	return nil, errors.New("not found")
}
