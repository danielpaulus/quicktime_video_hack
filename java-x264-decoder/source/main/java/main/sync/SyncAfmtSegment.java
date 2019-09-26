package main.sync;

import io.netty.buffer.ByteBuf;

public class SyncAfmtSegment extends SyncSegment {
	public SyncAfmtSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, header, type, data);
	}
}
