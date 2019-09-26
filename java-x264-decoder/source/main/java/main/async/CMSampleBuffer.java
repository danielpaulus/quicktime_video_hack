package main.async;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.buffer.Unpooled;
import main.NaluTypeTable;

public class CMSampleBuffer {

	static final int sbuf = 0x73627566;
	static final int opts = 0x6F707473;
	static final int stia = 0x73746961;
	static final int sdat = 0x73646174;

	private String options;
	private String stiaVal;
	private ByteBuf streamData;

	public CMSampleBuffer(ByteBuf data) {
		int length = data.readIntLE();
		int sbufmarker = data.readIntLE();
		if (sbufmarker != sbuf) {
			throw new IllegalStateException("oh nooo");
		}

		readOptions(data);
		readStia(data);
		readStreamData(data);

		//TODO: read metadata

	}

	private void readOptions(ByteBuf data) {
		int length = data.readIntLE();
		int optsMarker = data.readIntLE();
		if (optsMarker != opts) {
			throw new IllegalStateException("opts not found");
		}
		ByteBuf byteBuf = Unpooled.buffer(length - 8, length - 8);
		data.readBytes(byteBuf);
		options = ByteBufUtil.hexDump(byteBuf);
	}

	private void readStia(ByteBuf data) {
		int length = data.readIntLE();
		int stiaMarker = data.readIntLE();
		if (stiaMarker != stia) {
			throw new IllegalStateException("opts not found");
		}
		ByteBuf byteBuf = Unpooled.buffer(length - 8, length - 8);
		data.readBytes(byteBuf);
		stiaVal = ByteBufUtil.hexDump(byteBuf);
	}

	private void readStreamData(ByteBuf data) {
		int length = data.readIntLE();
		int stiaMarker = data.readIntLE();
		if (stiaMarker != sdat) {
			throw new IllegalStateException("opts not found");
		}
		ByteBuf byteBuf = Unpooled.buffer(length - 8, length - 8);
		data.readBytes(byteBuf);
		streamData = byteBuf;
	}

	@Override
	public String toString() {
		return "CMSampleBuffer{" +
				"options='" + options + '\'' +
				", stiaVal='" + stiaVal + '\'' +
				", streamData=" + streamData +
				", nalu:" + NaluTypeTable.getNaluDetails(streamData) +
				'}';
	}
}
