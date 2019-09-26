package main.async;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import main.RawSegment;

public class AsyncSegment extends RawSegment {

	final int type;
	ByteBuf header;

	public AsyncSegment(int length, int magicbytes, ByteBuf data, ByteBuf header, int type) {
		super(length, magicbytes, data);
		this.header = header;
		this.type = type;
	}

	public static AsyncSegment fromBytes(int length, int magicbytes, ByteBuf data) {
		ByteBuf header = Unpooled.buffer(8, 8);

		data.readBytes(header);
		int type = data.readIntLE();
		switch (type) {
		case AsyncSegmentTypes.FEED:
			return new AsyncFeedSegment(length, magicbytes, header, type, data);
		case AsyncSegmentTypes.SRAT:
			return new AsyncSratSegment(length, magicbytes, header, type, data);
		case AsyncSegmentTypes.TJMP:
			return new AsyncTJMPSegment(length, magicbytes, header, type, data);
		case AsyncSegmentTypes.TBAS:
			return new AsyncTBASegment(length, magicbytes, header, type, data);
		case AsyncSegmentTypes.SPRP:
			return new AsyncSPRPSegment(length, magicbytes, header, type, data);
		case AsyncSegmentTypes.RELS:
			return new AsyncRELSSegment(length, magicbytes, header, type, data);
		default:

			throw new RuntimeException("Unkown asyncsegment type" + ByteBufUtil.hexDump(data));
		}
	}

	@Override
	public String toString() {
		return String.format("length:%d magic:%s header:%s type:%s data:%s", length, toASCII(magicbytes), ByteBufUtil.hexDump(header),
				toASCII(type), ByteBufUtil.hexDump(data));

	}
}
