package main.dict;

import io.netty.buffer.ByteBuf;

import java.awt.*;
import java.nio.charset.Charset;
import java.util.LinkedList;
import java.util.List;

import static main.dict.DictTypes.extn;

public class FormatDescription {

	static final int vdim = 0x7664696D;

	private List<Object> properties = new LinkedList<>();

	private Dict extensions;

	public FormatDescription(ByteBuf pairData, int totalLength) {
		int initialIndex = pairData.readerIndex();
		while (pairData.readerIndex() - initialIndex < totalLength) {
			int length = pairData.readIntLE();
			parseVal(length - 4, pairData);
		}

	}

	private void parseVal(int length, ByteBuf pairData) {
		int marker = pairData.getIntLE(pairData.readerIndex());
		if (marker == vdim) {
			pairData.readIntLE(); //drop marker
			int x = pairData.readIntLE();
			int y = pairData.readIntLE();
			properties.add(new Dimension(x, y));
			return;
		}
		if (marker == extn) {
			extensions = new Dict(pairData);
			return;
		}
		CharSequence ascii = pairData.readCharSequence(length, Charset.forName("ascii"));
		properties.add(new StringBuilder(ascii).reverse());
	}

	@Override
	public String toString() {
		return "FormatDescription{" +
				"properties=" + properties +
				", extensions=" + extensions +
				'}';
	}
}
