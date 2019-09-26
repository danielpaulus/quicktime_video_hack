package main;

import io.netty.buffer.ByteBuf;

public class NaluTypeTable {
	//https://www.semanticscholar.org/paper/Multiplexing-the-elementary-streams-of-H.264-video-Siddaraju-Rao/c7b0e625198b663be9d61c3ec7e1ec341627168c/figure/0
	//for debugging purposes

	static final String[] naluTypes = new String[] { "unspecified", "coded slice", "data partition A",
			"data partition B", "data partition C", "IDR", "SEI", "sequence parameter set", "picture parameter set",
			"access unit delim", "end of seq", "end of stream", "filler data",
			"extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended", "extended",
			"undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined", "undefined" };

	public static String getNaluDetails(ByteBuf nalu) {
		return String.format("Nalu length:%d type:%s", nalu.getInt(nalu.readerIndex()), getType(nalu.getByte(nalu.readerIndex() + 4)));
	}

	private static String getType(byte anInt) {
		int combiner = 0x1f;
		int result = combiner & anInt;
		return naluTypes[result];
	}
}
