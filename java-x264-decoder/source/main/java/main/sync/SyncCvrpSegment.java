package main.sync;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import main.dict.Dict;

public class SyncCvrpSegment extends SyncSegment {

	private final Dict dict;
	private ByteBuf unknown;

	public SyncCvrpSegment(int length, int magicbytes, ByteBuf header, int type, ByteBuf data) {
		super(length, magicbytes, header, type, data);
		dict = parseMe(data);
	}

	private Dict parseMe(ByteBuf data) {
		unknown = Unpooled.buffer(16, 16);
		data.readBytes(unknown);
		int dictLength = data.readIntLE() - 4; //the 4 bytes length also count
		ByteBuf dictBytes = Unpooled.buffer(dictLength, dictLength);
		data.readBytes(dictBytes);
		return new Dict(dictBytes);

	}

	@Override
	public String toString() {
		return "Sync cvrp, dict:" + dict.toString();
	}
}
