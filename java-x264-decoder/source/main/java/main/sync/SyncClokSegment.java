package main.sync;

import io.netty.buffer.ByteBuf;

public class SyncClokSegment extends SyncSegment {
	public SyncClokSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, header, type, data);
	}
}
