[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CircleCI](https://circleci.com/gh/danielpaulus/quicktime_video_hack.svg?style=svg)](https://circleci.com/gh/danielpaulus/quicktime_video_hack)
[![codecov](https://codecov.io/gh/danielpaulus/quicktime_video_hack/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpaulus/quicktime_video_hack)
[![Go Report](https://goreportcard.com/badge/github.com/danielpaulus/quicktime_video_hack)](https://goreportcard.com/report/github.com/danielpaulus/quicktime_video_hack)

## 1. What is this?
This is an Operating System indepedent implementation for Quicktime Screensharing for iOS devices :-)

[Check out my presentation](https://danielpaulus.github.io/quicktime_video_hack_presentation)

[See a demo on YouTube](https://youtu.be/8v5f_ybSjHk)

This repository contains all the code you will need to grab and record video and audio from your iPhone or iPad 
without needing one of these expensive MacOS X computers :-D
It probably does something similar to what `QuickTime` and `com.apple.cmio.iOSScreenCaptureAssistant` are doing on MacOS.
Currently you can use it to create a h264 file and a wave file so you can watch and listen what was happening on your device. 
I am finishing up and RTP implementation so you can live watch and hear your device. 
If you want to contribute to code or documentation, please go ahead :-D

## 2. Technical Docs
I have written some documentation here [doc/technical_documentation.md](https://github.com/danielpaulus/quicktime_video_hack/blob/master/doc/technical_documentation.md)
So if you are just interested in the protocol or if you want to implement this in a different programming language than golang, read the docs.
## 3. Usage& Current State of the Tool
- run `qvh --help` to see how it works
- The `record` command lets you save iOS video and Audio into separate h264 and wave files.
- The `gstreamer` command will render the video in a nice window on your screen. For this to work you need to install gstreamer on you machine.

Progress:

0. ~Fix it for MacOS X~

0a. Release 0.1-beta

1. Use Gstreamer to create a nice Window with Video and Audio
2. Release 0.2-beta
3. Allow Creating MPEG files
4. Release 0.3-beta
5. BUG- Linux Only: After running the tool to grab AV data, you have to unplug the device to make it work again
6. BUG- MacOSX Only: After running the tool, you have to wait a while until it works again (or keep restarting the tool)

Extra Goals:

1. [Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)



## 4. Additional Notes
### MAC OS X LIBUSB -- IMPORTANT
1. Make sure to use either this fork `https://github.com/GroundControl-Solutions/libusb`
   or a LibUsb version BELOW 1.0.20 or iOS devices won't be found on Mac OS X.
   [See Github Issue](https://github.com/libusb/libusb/issues/290)

