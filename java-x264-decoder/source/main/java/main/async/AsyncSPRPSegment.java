package main.async;

import io.netty.buffer.ByteBuf;

public class AsyncSPRPSegment extends AsyncSegment {
	public AsyncSPRPSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
	}
}
