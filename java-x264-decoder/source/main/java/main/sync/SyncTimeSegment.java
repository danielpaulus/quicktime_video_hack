package main.sync;

import io.netty.buffer.ByteBuf;

public class SyncTimeSegment extends SyncSegment {
	public SyncTimeSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, header, type, data);
	}
}
