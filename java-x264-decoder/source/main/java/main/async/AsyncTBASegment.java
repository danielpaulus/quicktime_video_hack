package main.async;

import io.netty.buffer.ByteBuf;

public class AsyncTBASegment extends AsyncSegment {
	public AsyncTBASegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
	}
}
