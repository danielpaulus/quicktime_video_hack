[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/danielpaulus/quicktime_video_hack)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CircleCI](https://circleci.com/gh/danielpaulus/quicktime_video_hack.svg?style=svg)](https://circleci.com/gh/danielpaulus/quicktime_video_hack)
[![codecov](https://codecov.io/gh/danielpaulus/quicktime_video_hack/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpaulus/quicktime_video_hack)
[![Go Report](https://goreportcard.com/badge/github.com/danielpaulus/quicktime_video_hack)](https://goreportcard.com/report/github.com/danielpaulus/quicktime_video_hack)

Release 0.6

- qvh without Gstreamer is finally stable on MacOSX. I ran it for 16 hours straight on parallel devices and it worked flawlessly.
- before a 1.0 Release I need to see if Gstreamer is stable enough and maybe fix or switch to ffmpeg
- Linux support needs to be improved. It works but it is hard to get going currently.
- Create an issue if you miss anything

## 1. What is this?

This is an Operating System indepedent implementation for Quicktime Screensharing for iOS devices :-)

[Check out my presentation](https://danielpaulus.github.io/quicktime_video_hack_presentation)

[See me talk about it at GoWayFest](https://www.youtube.com/watch?v=jghi4nCBRwc)

[See a demo on YouTube](https://youtu.be/8v5f_ybSjHk)

This repository contains all the code you will need to grab and record video and audio from one or more iPhone(s) or iPad(s)
without needing one of these expensive MacOS X computers or the hard to use QuickTime Player :-D

- You can record video and audio as raw h264 and wave audio in the Apple demonstration mode (Device shows 9:41am, full battery and no cellphone carrier in the status bar)
- Also you can just grab device audio as wave, ogg or mp3 without the Apple demonstration mode now 🎉
- You can use custom Gstreamer Pipelines to transcode the AV data into whatever you like

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

- For just displaying the screen run `qvh gstreamer` and it will work.
- For just getting raw media output without Gstreamer involved use `qvh record out.h264 out.wav` or `qvh audio out.wav --wav` for audio only
- For troubleshooting run `qvh diagnostics metrics.csv --dump=binary.bin` which will persist logs to a file, dump all usb transfers and gather metrics.
- See `qvh gstreamer --examples` for transcoding media or streaming.
- For creating mp3 or ogg in audio only mode see `qvh audio out.mp3 --mp3` and `qvh audio out.ogg --ogg`

## 4. Technical Docs/ Roll your own implementation

QVH probably does something similar to what `QuickTime` and `com.apple.cmio.iOSScreenCaptureAssistant` are doing on MacOS.
I have written some documentation here [doc/technical_documentation.md](https://github.com/danielpaulus/quicktime_video_hack/blob/master/doc/technical_documentation.md)
So if you are just interested in the protocol or if you want to implement this in a different programming language than golang, read the docs.
Also I have extracted binary dumps of all messages for writing unit tests and re-develop this in your preferred language in a test driven style.

I have given up on windows support :-)
~~[Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)~~ Did not find a way to do it
