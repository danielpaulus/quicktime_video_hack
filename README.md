[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CircleCI](https://circleci.com/gh/danielpaulus/quicktime_video_hack.svg?style=svg)](https://circleci.com/gh/danielpaulus/quicktime_video_hack)
[![codecov](https://codecov.io/gh/danielpaulus/quicktime_video_hack/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpaulus/quicktime_video_hack)
[![Go Report](https://goreportcard.com/badge/github.com/danielpaulus/quicktime_video_hack)](https://goreportcard.com/report/github.com/danielpaulus/quicktime_video_hack)

## 1. What is this?
This is an Operating System indepedent implementation for Quicktime Screensharing for iOS devices :-)
This repository contains all the code you will need to grab and record video and audio from your iPhone or iPad 
without needing one of these expensive MacOS X computers :-D
It probably does something similar to what `QuickTime` and `com.apple.cmio.iOSScreenCaptureAssistant` are doing on MacOS.
Currently you can use it to create a h264 file that is playable with VLC and watch a recording of your devices screen :-D
If you want to contribute to code or documentation, please go ahead :-D

## 2. Technical Docs
I have written some documentation here [doc/technical_documentation.md](https://github.com/danielpaulus/quicktime_video_hack/blob/master/doc/technical_documentation.md)
So if you are just interested in the protocol or if you want to implement this in a different programming language than golang, read the docs.
## 3. Usage& Current State of the Tool
run `go run main.go --help` to see how it works

To run RTP you can use this gstreamer command 
```
gst-launch-1.0 -v udpsrc port=4000 buffer-size=622080 caps = "application/x-rtp, media=(string)video, clock-rate=(int)90000, encoding-name=(string)H264, payload=(int)96,framerate=(fraction)60/1" ! rtpjitterbuffer latency=100 ! rtph264depay  ! queue! decodebin ! queue ! videoconvert ! xvimagesink sync=false

```
and then run `go run main.go rtpstream localhost 4000`

Progress:
1. Send correct replies for clock SKEW packets
2. Fix small bug in lengthfield based decoder
3. Stream device audio over rtp as well
4. Make a release :-D


Extra Goals:

0. ~~Also save the device audio stream (I am already decoding it and receiving it, just not doing anything with it for now)~~
1. [Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)
2. See if there is maybe a way to get it to work on mac
3. BUG: After running the tool to grab AV data, you have to unplug the device to record another video currently

## 4. Additional Notes
### MAC OS X LIBUSB -- IMPORTANT
1. What works:
 You can enable the QuickTime config and discover QT capable devices with `qvh devices` and  `qvh activate` 

2. What does not work
Recording or streaming AV sessions won't work on MAC OS. I cannot claim the USB endpoint, do not know why currently. Maybe it is already claimed or I need to codesign my binary.  

2. Make sure to use either this fork `https://github.com/GroundControl-Solutions/libusb`
   or a LibUsb version BELOW 1.0.20 or iOS devices won't be found on Mac OS X.
   [See Github Issue](https://github.com/libusb/libusb/issues/290)

