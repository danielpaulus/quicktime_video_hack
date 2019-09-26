package main.async;

import io.netty.buffer.ByteBuf;

//Maybe time jump
public class AsyncTJMPSegment extends AsyncSegment {
	public AsyncTJMPSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
	}
}
