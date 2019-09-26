## Read only branch that contains the Java Decoder 
###  Operating System indepedent implementation for Quicktime Screensharing for iOS devices :-)
This repository contains all the code you will need to grab and record video and audio from your iPhone or iPad 
without needing one of these expensive MacOS X computers :-D
It probably does something similar to what `QuickTime` and `com.apple.cmio.iOSScreenCaptureAssistant` are doing on MacOS.

Progress:
1. ~~Write a prototype decoder in Java to reconstruct a Video from a USB Dump taken with Wireshark~~
2. ~~Create a Golang App to successfully enable and grab the first video stream data from the iOS Device on Linux~~
3. Port decoder from Java to Golang
4. Generate SPS and PPS dynamically
5. Generate GStreamer compatible x264 stream
6. Port to Windows (I don't know why, but still people use Windows nowadays)

run the `qvh` tool to get details :-)

### MAC OS X LIBUSB -- IMPORTANT
Make sure to use either this fork `https://github.com/GroundControl-Solutions/libusb`
or a LibUsb version BELOW 1.0.20 or iOS devices won't be found on Mac OS X.
[See Github Issue](https://github.com/libusb/libusb/issues/290)
