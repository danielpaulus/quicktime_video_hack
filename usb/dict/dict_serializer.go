package dict

import (
	"encoding/binary"
	"log"
)

func SerializeStringKeyDict(stringKeyDict StringKeyDict) []byte {
	buffer := make([]byte, 1024*1024)
	var slice = buffer[8:]
	var index = 0
	for _, entry := range stringKeyDict.Entries {
		keyvaluePair := slice[index+8:]
		keyLength := serializeKey(entry.Key, keyvaluePair)
		valueLength := serializeValue(entry.Value, keyvaluePair[keyLength:])
		writeLengthAndMagic(slice[index:], keyLength+valueLength+8, KeyValuePairMagic)
		index += 8 + valueLength + keyLength
	}
	dictSizePlusHeaderAndLength := index + 4 + 4
	writeLengthAndMagic(buffer, dictSizePlusHeaderAndLength, DictionaryMagic)

	return buffer[:dictSizePlusHeaderAndLength]
}

func serializeValue(value interface{}, bytes []byte) int {
	switch value.(type) {
	case bool:
		writeLengthAndMagic(bytes, 9, BooleanValueMagic)
		var boolValue uint32
		if value.(bool) {
			boolValue = 1
		}
		binary.LittleEndian.PutUint32(bytes[8:], boolValue)
		return 9
	case NSNumber:
		numberBytes := value.(NSNumber).ToBytes()
		length := len(numberBytes) + 8
		writeLengthAndMagic(bytes, length, NumberValueMagic)
		copy(bytes[8:], numberBytes)
		return length
	case string:
		stringValue := value.(string)
		length := len(stringValue) + 8
		writeLengthAndMagic(bytes, length, StringValueMagic)
		copy(bytes[8:], stringValue)
		return length
	case []byte:
		byteValue := value.([]byte)
		length := len(byteValue) + 8
		writeLengthAndMagic(bytes, length, DataValueMagic)
		copy(bytes[8:], byteValue)
		return length
	case StringKeyDict:
		dictValue := SerializeStringKeyDict(value.(StringKeyDict))
		copy(bytes, dictValue)
		return len(dictValue)
	default:
		log.Fatalf("Wrong type while serializing dict:%s", value)
	}
	return 0
}

func serializeKey(key string, bytes []byte) int {
	keyLength := len(key) + 8
	writeLengthAndMagic(bytes, keyLength, StringKey)
	copy(bytes[8:], key)
	return keyLength
}
