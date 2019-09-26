package main.sync;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import main.RawSegment;

public class SyncSegment extends RawSegment {

	protected final int type;
	protected ByteBuf header;

	public SyncSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, data);
		this.header = header;
		this.type = type;
	}

	public static SyncSegment fromBytes(int length, int magicbytes, ByteBuf data) {
		ByteBuf header = Unpooled.buffer(8, 8);

		data.readBytes(header);
		int type = data.readIntLE();
		switch (type) {
		case SyncSegmentTypes.TIME:
			return new SyncTimeSegment(length, magicbytes, header, type, data);
		case SyncSegmentTypes.AFMT:
			return new SyncAfmtSegment(length, magicbytes, header, type, data);
		case SyncSegmentTypes.CLOK:
			return new SyncClokSegment(length, magicbytes, header, type, data);
		case SyncSegmentTypes.CVRP:
			return new SyncCvrpSegment(length, magicbytes, header, type, data);
		case SyncSegmentTypes.CWPA:
			return new SyncCwpaSegment(length, magicbytes, header, type, data);

		default:
			throw new RuntimeException("Unkown syncsegment type" + ByteBufUtil.hexDump(data));
		}
	}

	@Override
	public String toString() {
		return String.format("length:%d magic:%s header:%s type:%s data:%s", length, toASCII(magicbytes), ByteBufUtil.hexDump(header),
				toASCII(type), ByteBufUtil.hexDump(data));

	}

}
