package main;

import java.nio.file.Files;
import java.nio.file.Paths;

public class Run {

	public static void main(String... args) throws Exception {
		byte[] bytes = Files.readAllBytes(Paths.get(
				"usbdump-videorecording/src-32.26.6_deviceaddress-26_endpointaddress-0x00000086_endpointaddressnumber-6.bin"));
		Parser p = new Parser();
		p.parse(bytes);

	}

}
