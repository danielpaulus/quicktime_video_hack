package main.async;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import main.RawSegment;

//could be release
public class AsyncRELSSegment extends AsyncSegment {
	public AsyncRELSSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
	}

	@Override
	public String toString() {
		return String.format("length:%d magic:%s header:%s type:%s has no data", length, RawSegment.toASCII(magicbytes),
				ByteBufUtil.hexDump(header),
				RawSegment.toASCII(type));

	}
}
