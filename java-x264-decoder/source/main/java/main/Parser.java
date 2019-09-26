package main;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import main.async.AsyncSegment;
import main.sync.SyncSegment;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
import java.util.List;

import static main.SegmentTypes.*;

public class Parser {

	private static final Logger logger = LoggerFactory.getLogger(Parser.class);

	public List<RawSegment> parse(byte[] bytes) {
		return parse(Unpooled.wrappedBuffer(bytes));
	}

	List<RawSegment> parse(ByteBuf wrappedBuffer) {
		var segments = new LinkedList<RawSegment>();
		while (wrappedBuffer.readableBytes() != 0) {
			RawSegment rawSegment = readSegment(wrappedBuffer);
			segments.add(rawSegment);
			logger.info("{}", rawSegment);
		}
		return segments;
	}

	private RawSegment readSegment(ByteBuf wrappedBuffer) {
		int length = wrappedBuffer.readIntLE();
		int magicbytes = wrappedBuffer.readIntLE();
		ByteBuf data = Unpooled.buffer(length - 8);
		wrappedBuffer.readBytes(data, length - 8);
		switch (magicbytes) {
		case PING:
			return new PingSegment(length, magicbytes, data);
		case SYNC:
			return SyncSegment.fromBytes(length, magicbytes, data);
		case ASYNC:
			return AsyncSegment.fromBytes(length, magicbytes, data);
		default:
			return new RawSegment(length, magicbytes, data);
		}

	}
}
