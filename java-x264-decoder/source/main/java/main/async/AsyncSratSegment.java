package main.async;

import io.netty.buffer.ByteBuf;

//maybe sync rate?
public class AsyncSratSegment extends AsyncSegment {
	public AsyncSratSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
	}
}
