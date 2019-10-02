package dict

import (
	"encoding/binary"
	"log"
	"os"
)

func SerializeStringKeyDict(stringKeyDict StringKeyDict) []byte {
	buffer := make([]byte, 1024*1024)
	var slice = buffer[8:]
	var index = 0
	for _, entry := range stringKeyDict.Entries {
		keyvaluePair := slice[8:]
		keyLength := serializeKey(entry.Key, keyvaluePair)
		valueLength := serializeValue(entry.Value, keyvaluePair[keyLength:])
		writeLengthAndMagic(slice, keyLength+valueLength+8, KeyValuePairMagic)
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
	}
	log.Fatal("wrong type")
	os.Exit(1)
	return 0
}

func serializeKey(key string, bytes []byte) int {
	keyLength := len(key) + 8
	writeLengthAndMagic(bytes, keyLength, StringKey)
	copy(bytes[8:], key)
	return keyLength
}
