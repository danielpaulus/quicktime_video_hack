package h264

import (
	"encoding/binary"
	"fmt"
)

//https://yumichan.net/video-processing/video-compression/introduction-to-h264-nal-unit/
//https://www.semanticscholar.org/paper/Multiplexing-the-elementary-streams-of-H.264-video-Siddaraju-Rao/c7b0e625198b663be9d61c3ec7e1ec341627168c/figure/0
//for debugging purposes

var naluTypes = Table()

func Table() []string {
	return []string{"unspecified", "coded slice", "data partition A",
		"data partition B", "data partition C", "IDR", "SEI", "sequence parameter set", "picture parameter set",
		"access unit delim", "end of seq", "end of stream", "filler data",
		"extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended",
		"undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined"}
}

func GetNaluDetails(nalu []byte) string {
	return fmt.Sprintf("Nalu length:%d type:%s", binary.BigEndian.Uint32(nalu), getType(nalu[4]))
}

func getType(anInt byte) string {
	combiner := 0x1f
	result := combiner & int(anInt)
	return naluTypes[result]
}
