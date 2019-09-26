package main.sync;

import io.netty.buffer.ByteBuf;

public class SyncCwpaSegment extends SyncSegment {
	public SyncCwpaSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, header, type, data);
	}
}
