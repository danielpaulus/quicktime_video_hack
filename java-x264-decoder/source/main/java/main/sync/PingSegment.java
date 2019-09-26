package main;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;

public class PingSegment extends RawSegment {
	public PingSegment(int length, int magicbytes, ByteBuf data) {
		super(length, magicbytes, data);
	}

	@Override
	public String toString() {
		return String.format("length:%d magic:%s data:%s", length, toASCII(magicbytes), ByteBufUtil.hexDump(data));
	}

}
