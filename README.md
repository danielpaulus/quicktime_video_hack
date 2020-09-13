[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/danielpaulus/quicktime_video_hack) 

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
Currently you can use it to dump a h264 file and a wave file or mirror your device in a desktop window. Transcoding it to anything else is very easy since I used Gstreamer to render the AV data. 

## 2. Installation
### 2.1 Mac OSX
1. On MacOS run `brew install libusb pkg-config gstreamer gst-plugins-bad gst-plugins-good gst-plugins-base gst-plugins-ugly`
2. To just run: Download the latest release and run it
3. To develop: Clone the repo and execute `go run main.go` (need to install golang of course)
### 2.2 Linux
1. Run with Docker: the Docker files are [here](https://github.com/danielpaulus/quicktime_video_hack/tree/master/docker). There is one for just building and one for running. 

2. If you want to build/run locally then copy paste the dependencies from this [Dockerfile](https://github.com/danielpaulus/quicktime_video_hack/blob/master/docker/Dockerfile.debian) and install with apt.
3. Git clone the repo and start hacking or download the latest release and run the binary :-D


## 3. Usage
```
Q.uickTime V.ideo H.ack (qvh) v0.5-beta

Usage:
  qvh devices [-v]
  qvh activate [--udid=<udid>] [-v]
  qvh record <h264file> <wavfile> [--udid=<udid>] [-v]
  qvh audio <outfile> (--mp3 | --ogg | --wav) [--udid=<udid>] [-v]
  qvh gstreamer [--pipeline=<pipeline>] [--examples] [--udid=<udid>] [-v]
  qvh --version | version


Options:
  -h --help       Show this screen.
  -v              Enable verbose mode (debug logging).
  --version       Show version.
  --udid=<udid>   UDID of the device. If not specified, the first found device will be used automatically.

The commands work as following:
	devices		lists iOS devices attached to this host and tells you if video streaming was activated for them

	activate	enables the video streaming config for the device specified by --udid

	record		will start video&audio recording. Video will be saved in a raw h264 file playable by VLC.
	            	Audio will be saved in a uncompressed wav file. Run like: "qvh record /home/yourname/out.h264 /home/yourname/out.wav"

	audio		Records only audio from the device. It does not change the status bar like the video recording mode does.
			The recorded audio will be saved in <outfile> with the selected format. Currently (--mp3 | --ogg | --wav) are supported.
			Adding more formats is trivial though so create an issue or a PR if you need something :-)

	gstreamer	If no additional param is provided, qvh will open a new window and push AV data to gstreamer.
			If "qvh gstreamer --examples" is provided, qvh will print some common gstreamer pipeline examples.
			If --pipeline is provided, qvh will use the provided gstreamer pipeline instead of
			displaying audio and video in a window.
```

## 3. Technical Docs/ Roll your own implementation
I have written some documentation here [doc/technical_documentation.md](https://github.com/danielpaulus/quicktime_video_hack/blob/master/doc/technical_documentation.md)
So if you are just interested in the protocol or if you want to implement this in a different programming language than golang, read the docs.
Also I have extracted binary dumps of all messages for writing unit tests and re-develop this in your preferred language in a test driven style.

I have given up on windows support  :-)
~~[Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)~~ Did not find a way to do it



