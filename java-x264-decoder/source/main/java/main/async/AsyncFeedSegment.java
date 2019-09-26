package main.async;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import main.RawSegment;

public class AsyncFeedSegment extends AsyncSegment {

	private final CMSampleBuffer cmSampleBuffer;

	public AsyncFeedSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data, header, type);
		cmSampleBuffer = new CMSampleBuffer(data);
	}

	@Override
	public String toString() {
		return String
				.format("length:%d magic:%s header:%s type:%s \n %s", length, RawSegment.toASCII(magicbytes), ByteBufUtil.hexDump(header),
						RawSegment.toASCII(type), cmSampleBuffer);
	}
}
