package main.dict;

import com.google.common.base.Joiner;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;

import java.nio.charset.Charset;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Collectors;

import static main.RawSegment.toASCII;

public class Dict {

	private List<DictPair> entries = new LinkedList<>();

	public Dict(ByteBuf dictBytes) {
		int dictMarker = dictBytes.readIntLE();

		while (dictBytes.readableBytes() > 0) {
			int entryLength = dictBytes.readIntLE();
			int entryMarker = dictBytes.readIntLE();

			ByteBuf pairData = Unpooled.buffer(entryLength - 8, entryLength - 8);
			dictBytes.readBytes(pairData);
			Object key = readKey(pairData);
			Object entry = readEntry(pairData);
			DictPair pair = new DictPair(key.toString(), entry);
			entries.add(pair);
		}
	}

	private Object readEntry(ByteBuf pairData) {
		int length = pairData.readIntLE();
		int type = pairData.getIntLE(pairData.readerIndex());
		String s = toASCII(type);
		if (type == DictTypes.dict) {
			return new Dict(pairData);
		}
		pairData.readInt();
		if (type == DictTypes.nmbv) {
			ByteBuf byteBuf = Unpooled.buffer(length - 8, length - 8);
			pairData.readBytes(byteBuf);
			return ByteBufUtil.hexDump(byteBuf);
		}

		if (type == DictTypes.fdsc) {
			return new FormatDescription(pairData, length - 8);
		}

		if (type == DictTypes.datv) {
			ByteBuf byteBuf = Unpooled.buffer(length - 8, length - 8);
			pairData.readBytes(byteBuf);
			return ByteBufUtil.hexDump(byteBuf);
		}

		return pairData.readCharSequence(length - 8, Charset.forName("ascii"));
	}

	private Object readKey(ByteBuf pairData) {
		int length = pairData.readIntLE();
		int type = pairData.readIntLE();
		if (type == DictTypes.strk) {
			return pairData.readCharSequence(length - 8, Charset.forName("ascii"));
		}
		if (type == DictTypes.idxk) {
			return pairData.readShortLE();
		}
		throw new RuntimeException("unknown dict key type:" + type);
	}

	@Override
	public String toString() {
		List<String> collect = entries.stream().map(DictPair::toString).collect(Collectors.toList());
		return "[" + Joiner.on(",").join(collect) + "]";
	}
}
