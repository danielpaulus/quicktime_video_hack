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

Progress:
1. ~~Make the `go run main.go record` work on the first execution (currently you have to run it twice and it will start recording on the second run)~~
2. FIX: After running the record command and saving a video, you have to unplug the device to record another video currently
3. Make a release :-D
4. ~~Generate GStreamer compatible x264 stream probably by wrapping the NaLus in RTP headers~~
5. ~~Complete packet documentation~~
6. Send correct replies for clock SKEW packets
7. Fix small bug in lengthfield based decoder
8. Stream device audio over rtp as well


Extra Goals:

0. ~~Also save the device audio stream (I am already decoding it and receiving it, just not doing anything with it for now)~~
1. [Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)


## 4. Additional Notes
### MAC OS X LIBUSB -- IMPORTANT
1. What works:
 You can enable the QuickTime config and discover QT capable devices with `qvh devices` and  `qvh activate` 

2. What does not work

This might be wrong, needs investigation--> `qvh record` won't work on MAC OS because the binary needs to be codesigned with `com.apple.ibridge.control`
 apparently that is a protected Entitlement that I have no idea how to use or sign my binary with. 

2. Make sure to use either this fork `https://github.com/GroundControl-Solutions/libusb`
   or a LibUsb version BELOW 1.0.20 or iOS devices won't be found on Mac OS X.
   [See Github Issue](https://github.com/libusb/libusb/issues/290)

