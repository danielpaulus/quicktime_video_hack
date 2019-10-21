package coremedia

import (
	"encoding/binary"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/common"
	"log"
)

//SerializeStringKeyDict serializes a StringKeyDict into a []byte
func SerializeStringKeyDict(stringKeyDict StringKeyDict) []byte {
	buffer := make([]byte, 1024*1024)
	var slice = buffer[8:]
	var index = 0
	for _, entry := range stringKeyDict.Entries {
		keyvaluePair := slice[index+8:]
		keyLength := serializeKey(entry.Key, keyvaluePair)
		valueLength := serializeValue(entry.Value, keyvaluePair[keyLength:])
		common.WriteLengthAndMagic(slice[index:], keyLength+valueLength+8, KeyValuePairMagic)
		index += 8 + valueLength + keyLength
	}
	dictSizePlusHeaderAndLength := index + 4 + 4
	common.WriteLengthAndMagic(buffer, dictSizePlusHeaderAndLength, DictionaryMagic)

	return buffer[:dictSizePlusHeaderAndLength]
}

func serializeValue(value interface{}, bytes []byte) int {
	switch value := value.(type) {
	case bool:
		common.WriteLengthAndMagic(bytes, 9, BooleanValueMagic)
		var boolValue uint32
		if value {
			boolValue = 1
		}
		binary.LittleEndian.PutUint32(bytes[8:], boolValue)
		return 9
	case common.NSNumber:
		numberBytes := value.ToBytes()
		length := len(numberBytes) + 8
		common.WriteLengthAndMagic(bytes, length, common.NumberValueMagic)
		copy(bytes[8:], numberBytes)
		return length
	case string:
		stringValue := value
		length := len(stringValue) + 8
		common.WriteLengthAndMagic(bytes, length, StringValueMagic)
		copy(bytes[8:], stringValue)
		return length
	case []byte:
		byteValue := value
		length := len(byteValue) + 8
		common.WriteLengthAndMagic(bytes, length, DataValueMagic)
		copy(bytes[8:], byteValue)
		return length
	case StringKeyDict:
		dictValue := SerializeStringKeyDict(value)
		copy(bytes, dictValue)
		return len(dictValue)
	default:
		log.Fatalf("Wrong type while serializing dict:%s", value)
	}
	return 0
}

func serializeKey(key string, bytes []byte) int {
	keyLength := len(key) + 8
	common.WriteLengthAndMagic(bytes, keyLength, StringKey)
	copy(bytes[8:], key)
	return keyLength
}
