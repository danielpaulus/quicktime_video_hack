package main;

import io.netty.buffer.ByteBuf;

public class RawSegment {

	final static Parser PARSER = new Parser();

	protected final int length;
	protected final int magicbytes;
	protected final ByteBuf data;

	public RawSegment(int length, int magicbytes, ByteBuf data) {
		this.length = length;
		this.magicbytes = magicbytes;
		this.data = data;
	}

	public static String toASCII(int value) {
		int length = 4;
		StringBuilder builder = new StringBuilder(length);
		for (int i = length - 1; i >= 0; i--) {
			builder.append((char) ((value >> (8 * i)) & 0xFF));
		}
		return builder.toString();
	}

	@Override
	public String toString() {
		return String.format("length:%d magic:%s", length, toASCII(magicbytes));
	}

}
