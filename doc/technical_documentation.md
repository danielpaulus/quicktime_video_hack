# Table of Contents
- [Technical Documentation of the iOS ScreenSharing Feature](#technical-documentation-of-the-ios-screensharing-feature)
  * [0. General Information](#0-general-information)
  * [1. How to Enable it for a iOS Device on the USB Level](#1-how-to-enable-it-for-a-ios-device-on-the-usb-level)
    + [1.1 Foundations](#11-foundations)
    + [1.2 Finding Devices and Configurations using LibUsb](#12-finding-devices-and-configurations-using-libusb)
    + [1.3 Hidden Configuration](#13-hidden-configuration)
    + [1.4 Enabling the Hidden Config](#14-enabling-the-hidden-config)
  * [2. AV Session LifeCycle](#2-av-session-lifecycle)
    + [2.1 Initiate the session](#21-initiate-the-session)
    + [2.2 Receive data](#22-receive-data)
    + [2.3 Shutting down streaming](#23-shutting-down-streaming)
  * [3. Protocol Reference](#3-protocol-reference)
    + [3.1 Ping Packet](#31-ping-packet)
    + [3.2 Sync Packets](#32-sync-packets)
      - [3.2.1 General Description](#321-general-description)
      - [3.2.2. CWPA Packet and Response](#322-cwpa-packet-and-response)
        * [General Description](#general-description)
        * [Request Format Description](#request-format-description)
        * [Reply - RPLY Format Description](#reply---rply-format-description)
      - [3.2.3. AFMT Packet](#323-afmt-packet)
        * [General Description](#general-description-1)
        * [Request Format Description](#request-format-description-1)
        * [Reply - RPLY Format Description](#reply---rply-format-description-1)
      - [3.2.4. CVRP Packet](#324-cvrp-packet)
        * [General Description](#general-description-2)
        * [Request Format Description](#request-format-description-2)
        * [Reply - RPLY Format Description](#reply---rply-format-description-2)
      - [3.2.5. CLOK Packet](#325-clok-packet)
        * [General Description](#general-description-3)
        * [Request Format Description](#request-format-description-3)
        * [Reply - RPLY Format Description](#reply---rply-format-description-3)
      - [3.2.6. TIME Packet](#326-time-packet)
        * [General Description](#general-description-4)
        * [Request Format Description](#request-format-description-4)
        * [Reply - RPLY Format Description](#reply---rply-format-description-4)
      - [3.2.7. SKEW Packet](#327-skew-packet)
        * [General Description](#general-description-5)
        * [Request Format Description](#request-format-description-5)
        * [Reply - RPLY Format Description](#reply---rply-format-description-5)
      - [3.2.8. OG Packet](#328-og-packet)
        * [General Description](#general-description-6)
        * [Request Format Description](#request-format-description-6)
        * [Reply - RPLY Format Description](#reply---rply-format-description-6)
      - [3.2.9. STOP Packet](#329-stop-packet)
        * [General Description](#general-description-7)
        * [Request Format Description](#request-format-description-7)
        * [Reply - RPLY Format Description](#reply---rply-format-description-7)
    + [3.3 Asyn Packets](#33-asyn-packets)
        * [3.3.0 General Description](#330-general-description)
        * [3.3.1. Asyn SPRP - Set Property](#331-asyn-sprp---set-property)
          + [General Description](#general-description-8)
          + [Packet Format Description](#packet-format-description)
        * [3.3.2. Asyn SRAT - Set time rate and Anchor](#332-asyn-srat---set-time-rate-and-anchor)
          + [General Description](#general-description-9)
          + [Packet Format Description](#packet-format-description-1)
        * [3.3.3. Asyn TBAS - Set TimeBase](#333-asyn-tbas---set-timebase)
          + [General Description](#general-description-10)
          + [Packet Format Description](#packet-format-description-2)
        * [3.3.4. Asyn TJMP - Time Jump Notification](#334-asyn-tjmp---time-jump-notification)
          + [General Description](#general-description-11)
          + [Packet Format Description](#packet-format-description-3)
        * [3.3.5 Asyn FEED - CMSampleBuffer with h264 Video Data](#335-asyn-feed---cmsamplebuffer-with-h264-video-data)
          + [Packet Format Description](#packet-format-description-4)
        * [3.3.6 Asyn EAT! - CMSampleBuffer with Audio Data](#336-asyn-eat----cmsamplebuffer-with-audio-data)
        * [3.3.7 Asyn NEED - Tell the device to send more](#337-asyn-need---tell-the-device-to-send-more)
          + [Packet Format Description](#packet-format-description-5)
        * [3.3.8 Asyn HPD0 - Tell the device to stop video streaming](#338-asyn-hpd0---tell-the-device-to-stop-video-streaming)
          + [Packet Format Description](#packet-format-description-6)
        * [3.3.9 Asyn HPA0 - Tell the device to stop audio streaming](#339-asyn-hpa0---tell-the-device-to-stop-audio-streaming)
          + [Packet Format Description](#packet-format-description-7)
        * [3.3.10 Asyn RELS - Tell us about a released Clock on the device](#3310-asyn-rels---tell-us-about-a-released-clock-on-the-device)
          + [Packet Format Description](#packet-format-description-8)
  * [4 Serializing/Deserializing Objects](#4-serializing-deserializing-objects)
    + [4.0 General Description](#40-general-description)
    + [4.1 Dictionaries](#41-dictionaries)
      - [4.1.0 General Description](#410-general-description)
      - [4.1.1 Dictionaries with String Keys](#411-dictionaries-with-string-keys)
        * [4.1.1.1 General Dictionary Structure](#4111-general-dictionary-structure)
      - [4.1.2 Dictionaries with 4-Byte Integer Index Keys](#412-dictionaries-with-4-byte-integer-index-keys)
    + [4.2 CMTime](#42-cmtime)
      - [Example](#example)
    + [4.3 CMSampleBuffer](#43-cmsamplebuffer)
    + [4.4 NSNumber](#44-nsnumber)
      - [Example](#example-1)
    + [4.5 CMFormatDescription](#45-cmformatdescription)
  * [5. Clocks and CMSync](#5-clocks-and-cmsync)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Technical Documentation of the iOS ScreenSharing Feature
## 0. General Information
This document provides you with details about the screen sharing feature of QuickTime for iOS devices. 
The information contained in this document can be used to re-implement that feature in the programming language of choice and use
the feature on other operating systems than MAC OS X. If you want to implement the feature, I recommend using my unit tests and test fixtures. I have prepared a bin dump with an example for every message type. This way you can easily make sure your codec does what it's supposed to.
The repository also contains a reference implementation in Golang. 

-- Note: All the information in this document is reverse engineered by the author, therefore it could be wrong or not entirely accurate
as it involves a lot of assumptions. If you find mistakes or more accurate descriptions please add them :-) -- 


## 1. How to Enable it for a iOS Device on the USB Level
### 1.1 Foundations
Usually devices attached on the USB Port have a set of "configurations" that you can retrieve using the LibUsb wrapper you use.
Inside of these interfaces there are a set of Usb Endpoints you can use to communicate with your device. We are interested in the 
`bulk` endpoints of iOS devices as these are used for transferring data. 
### 1.2 Finding Devices and Configurations using LibUsb
Look for a USB ConfigDescriptor that has an interface with Class `DeviceClass = 0xFF` 
By default, iOS devices only have the USBMux ConfigDescriptor which has SubClass `0xFE` and will contain one interface with two Bulk endpoints.
The Subclass for a config that allows AV streaming will be  `0x2A` and there will be an interface with 4 Bulk endpoints. It needs to be enabled first.

### 1.3 Hidden Configuration
To use video mirroring, you have to enable the hidden Quicktime Configuration (I just call it QT-Config, not an official apple name)
If you closely monitor all available USB configurations you will find that if a device has n-Configurations, as soon as you open QuickTime on a mac
and start recording the iOS device's screen, there will be n+Configurations
### 1.4 Enabling the Hidden Config
To enable the hidden QTconfig you have to send a specific Control Request to the device like so:
`val, err := device.usbDevice.Control(0x40, 0x52, 0x00, 0x02, response)`
If you did it correctly, it will cause the device to disconnect from the host machine and re-connect after a few moments with an additional config.
The new config contains 4 bulk endpoints. 2 For communication with the usbmuxd on the device, and two additional endpoints for receiving and sending AV data.
Call setActiveConfiguration on that config and you can claim the new endpoint for sending and receiving AV data.
## 2. AV Session LifeCycle

### 2.1 Initiate the session
1. enable hidden device config
2. claim endpoint
3. wait to receive a PING packet
4. respond with a PING packet
5. wait for SYNC CWPA packet to receive the clockref for the devices audio clock
6. create local clock, put clockref in reply to the SYNC CWPA and send
7. send ASYN_HPD1
8. send ASYN_HPA1 with the device audio clockref received in step 6
9. receive SYNC AFMT and reply with a zero error code
10. receive SYNC CVRP with the devices video clockRef
11. reply with the local video clockRef
12. start sending ASYN NEED with the device's video clockRef
13. receive two ASYN Set Properties
14. receive Sync Clok and reply with newly created clock
15. receive two SYNC TIME and reply with two CMTimes 
### 2.2 Receive data
FEED and EAT! Packets for video and audio will be sent by the device.
We need to send NEED packets for video periodically

### 2.3 Shutting down streaming
1. send asyn hpa0 with the deviceclockref from the cwpa sync packet to tell the device to stop sending audio
2. send hpd0 with empty clockRef to stop video
3. receive sync stop package for our video clock we created when cvrp was sent to us and that is in every feed packet
4. reply to sync stop with 8 zero bytes
5. receive a ASYN RELS for the local video clockRef (the one found in FEED packets)
6. receive a ASYN RELS for the local clock created after the SYNC CLOCK 
7. release usb endpoint
8. set the device active config to usbmux only


## 3. Protocol Reference
### 3.1 Ping Packet
As soon as we connect to the USB endpoint, we need to wait for the device to send us a ping packet. Once we received it, we will send a ping back to the device and 
then progress to the rest of the communication. Example Ping:

| 4 Byte Length (16)   |4 Byte Magic (PING)   | 4byte 0x0| 4 byte 0x1   | 
|---|---|---|---|
|10000000 |676E6970 | 00000000 | 01000000 |
  

### 3.2 Sync Packets
#### 3.2.1 General Description
All SYNC packets require us to reply with a RPLY packet. 
It seems like this is mostly used for synchronizing CMClocks and exchanging 8byte CMClockRefs (implement CMSync.h protocol)
Usually you can see that SYNC packets have a 4byte SUB-TYPE followed by a 8byte correlationID. A reply always contains the correlationID so 
I assume this is how the device knows which reply belongs to which request.


#### 3.2.2. CWPA Packet and Response
##### General Description
This packet seems to be used for intitiating the audio stream. We get a clockRef from the device and respond with our own, newly created clockRef.
The clockref send by the device needs to go in the ASYN-1APH packet we send.
##### Request Format Description

| 4 Byte Length (36)   |4 Byte Magic (SYNC)   | 8 Empty clock reference| 4 byte message type (CWPA)   | 8 byte correlation id  | 8 bytes CFTypeID of the device clock |
|---|---|---|---|---|---|
|24000000 |636E7973 |01000000 00000000 | 61707763 |E03D5713 01000000| E0740000 5A130040 |

##### Reply - RPLY Format Description

Sends back our clockRef. The device will use the clockRef from here in the SYNC_AFMT message to tell us about the audio format. 
Also this will be used for all ASYN_EAT packets containing audio sample buffers. 

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |   4 Byte: 0  | 8 bytes CFTypeID of our clock |
|---|---|---|---|---|
|1C000000 | 796C7072 |E03D5713 01000000 | 00000000 |B00CE26C A67F0000|

#### 3.2.3. AFMT Packet
##### General Description
This packet contains information about the Audio Format(AMFT). It contains a AudioStreamBasicDescription struct [see here](https://github.com/nu774/MSResampler/blob/master/CoreAudio/CoreAudioTypes.h). It is usually of MediaID Linear pulse-code modulation (LPCM), which means uncompressed audio.
The response is basically a dictionary containing an error code. Normally we send 0 to indicate everything is ok.
Note how the device references the Clock we gave it in the SYNC_CWPA_RPLY
##### Request Format Description

| 4 Byte Length (68)   |4 Byte Magic (SYNC)   | 8 bytes clock CFTypeID| 4 byte magic (AFMT)| 8 byte correlation id| AudioStreamBasicDescription struct: float64 sampling frequency (48khz), data 4 byte magic (LPCM),   28 bytes rest|
|---|---|---|---|---|---|
|44000000| 636E7973| B00CE26C A67F0000| 746D6661 | 809D2213 01000000| 00000000 0070E740 6D63706C 4C000000 04000000 01000000 04000000 02000000 10000000 00000000|

##### Reply - RPLY Format Description
Contains the correlationID from the request as well as a simple Dictionary:  {"Error":NSNumberUint32(0)}

| 4 Byte Length (62)   |4 Byte Magic (RPLY)   | 8  correlation id| 4 byte 0| 4 byte dict length(42)| 4 byte magic (DICT)| dict bytes |
|---|---|---|---|---|---|---|
|3E000000 |796C7072| 809D2213 01000000 |00000000| 2A000000| 74636964| 22000000 7679656B 0D000000 6B727473 4572726F 720D0000 0076626D 6E030000 0000|

#### 3.2.4. CVRP Packet
##### General Description
The CVRP packet is used to create a local clock for synchronizing video frames. The ClockRef we get from the device, is the one we need to put in NEED packets.
Similar to this the device sends all ASYN_FEED sample buffer with a reference to our clock. After this, two SetProperty Asyns are usually received plus CLok, TBAS and SRAT.
##### Request Format Description

Contains a Dict with a FormatDescription and timing information. With the FormatDescription for the first time we get the h264 Picture and Sequence ParameterSet(PPS/SPS) already encoded in nice NALUs ready for streaming
over the network. They are hidden inside a dictionary inside the extension dictionary.

|4 Byte Length (649)|4 Byte Magic (SYNC)|8 byte empty(?) clock reference|4 byte magic(CVRP)|8 byte correlation id|CFTypeID of clock on device (needs to be in NEED packets we send)|4 byte length of dictionary (613)|4 byte magic (DICT)| Dict bytes|
|---|---|---|---|---|---|---|---|---|
|89020000 |636E7973| 01000000 00000000 |70727663| D0595613 01000000 |A08D5313 01000000 |65020000| 74636964|   0x.....|

##### Reply - RPLY Format Description

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |  4 bytes (seem to be always 0) | 8 bytes CFTypeID of our clock(will be in all feed async packets) |
|---|---|---|---|---|
|1C000000 | 796C7072 |D0595613 01000000 | 00000000 |5002D16C A67F0000 |


#### 3.2.5. CLOK Packet
##### General Description
This requests us to create a new clock and send back the clockRef. It is usually followed by 2 TIME requests.

##### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (CLOK) | 8 bytes correlation id |
|---|---|---|---|---|
|1C000000| 636E7973| 5002D16C A67F0000| 6B6F6C63 | 70495813 01000000 |

##### Reply - RPLY Format Description

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 correlation id  |  4 bytes (seem to be always 0) | 8 bytes CFTypeID of our clock(for the next two time packets) |
|---|---|---|---|---|
|1C000000| 796C7072| 70495813 01000000| 00000000 | 8079C17C A67F0000|

#### 3.2.6. TIME Packet
##### General Description
This packet requests from us to send a RPLY with the current CMTime for the ClockRef specified.
##### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (TIME) | 8 bytes correlation id |
|---|---|---|---|---|
|1C000000| 636E7973| 8079C17C A67F0000 |656D6974 | 503D2213 01000000 |

##### Reply - RPLY Format Description

| 4 Byte Length (44)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  |  4 bytes 0x0 | 24 bytes CMTime struct |
|---|---|---|---|---|
|2C000000 |796C7072 |503D2213 01000000| 00000000 | E1E142C4 62BA0000 00CA9A3B 01000000 00000000 00000000|

#### 3.2.7. SKEW Packet
##### General Description
This packet tells the device about the clock skew of the audio clock (clockRef used in EAT! packets, which we sent as response to cwpa). As denoted in this [wikipedia](https://en.wikipedia.org/wiki/Clock_skew#On_a_network) article, clock skew means the difference in frequency of both clocks. In other words, both clocks supposedly 
run at 48khz, and the device wants to know how many ticks per second our clock executed during the time the device clock had one tick. 
So we have to respond with:
- 48000 if the clocks were aligned
- some value above 48000 if our clock was slower
- and some value below 48000 if our clock was faster than the device clock
If implemented correctly, we should see that the skew responses converge towards 48000 with small deviations sometimes `(48000+x where -1 < x <1)`

##### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (SKEW) | 8 bytes correlation id | 
|---|---|---|---|---|
|20000000| 636E7973| 8079C17C A67F0000 |77656B73 | 60B9FD02 01000000 | 

##### Reply - RPLY Format Description

| 4 Byte Length (24)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  | 4bytes padding 0x0| 8 bytes floating point number (48000.0) | 
|---|---|---|---|---|
|18000000 |796C7072 |60B9FD02 01000000| 00000000 | 00000000 0070E740 |

#### 3.2.8. OG Packet
##### General Description
I do not know what this does or what it is for. It seems like it sends one uint32 as payload that is always equal to 1 and 
we have to reply back with a 8bytes zero
##### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (OG!.) | 8 bytes correlation id | 4byte unknown int|
|---|---|---|---|---|---|
|20000000| 636E7973| 8079C17C A67F0000 |20216F67 | 302FD302 01000000 | 01000000 |

##### Reply - RPLY Format Description

| 4 Byte Length (24)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  |  8 bytes 0x0 | 
|---|---|---|---|
|18000000 |796C7072 |302FD302 01000000| 00000000 00000000|

#### 3.2.9. STOP Packet
##### General Description
This one tells us to stop our clock
##### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (STOP) | 8 bytes correlation id | 
|---|---|---|---|---|
|1C000000| 636E7973| F05F4235 BA7F0000 | 706F7473 | 1049FD02 01000000 | 


##### Reply - RPLY Format Description

| 4 Byte Length (24)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  |  4 bytes 0x0 | 4 bytes 0x0 |
|---|---|---|---|---|
|18000000 |796C7072 |1049FD02 01000000| 00000000 | 00000000|

### 3.3 Asyn Packets
##### 3.3.0 General Description
Asyn packets contain information like CMFormatDescriptions, Properties and most importantly the CMSampleBuffers that contain audio and video data.
They start with 4 byte length, 4 byte magic and then a 8 byte ClockRef.
ASYN packets do not require us to respond or ack.
##### 3.3.1. Asyn SPRP - Set Property
###### General Description
This packet is used to set properties for the video stream. Usually you get only two of them containing a pair each referencing the video stream.
1. ObeyEmptyMediaMarkers = true
2. RenderEmptyMedia = false

###### Packet Format Description

| 4 Byte Length (varies)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (SPRP) | Key value pair bytes|
|---|---|---|---|---|
|20000000| 6E797361| 18BC2311 01000000 |70727073 | 00..... |


##### 3.3.2. Asyn SRAT - Set time rate and Anchor
###### General Description
I think this one along with TJMP, TBAS and CLOK is part of some kind of clock synchronization algorithm. 
I believe it is related to what AVPlayer.SetRate usually does. 
It could also tell us to invoke this on our clock: CMTimebaseSetRateAndAnchorTime https://developer.apple.com/documentation/coremedia/cmtimebase?language=objc

###### Packet Format Description


| 4 Byte Length (49)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (SRAT) | 32bit float (1.0)| 32bit float (1.0)| 24 byte CMTime|
|---|---|---|---|---|---|---|
|31000000| 6E797361| 18BC2311 01000000 |74617273 | 0000803F | 0000803F| CB448EA1 CF10CC15 00CA9A3B 01000000 00000000 00000000|

##### 3.3.3. Asyn TBAS - Set TimeBase
###### General Description
This one contains a clockRef and another clockRef that I do not know about. It could be that we are supposed to create a CMTimeBase for our clock with the 
given ref. Also I think it could be that the device created a CMTimeBase and just tells us about it. 
So far, I could not find another usage of the Unknown Clock/TimeBaseRef we get in the whole communication. 

###### Packet Format Description

| 4 Byte Length (24)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (TBAS) | 8byte Unknown ClockRef|
|---|---|---|---|---|
|18000000| 6E797361| 18BC2311 01000000 |73616274 | C0904402 01000000|

##### 3.3.4. Asyn TJMP - Time Jump Notification
###### General Description
I think this packet tells us that a CMTimeBase on the device was set to a different time. 
As this is not referenced or used anywhere else, I don't know exactly what it means or what the values are for. 

The payload is 56 bytes. I think it could be like:  4byte int 0x0, 4byte 0x0, CMTime, CMTime but i do not know for sure.
###### Packet Format Description

| 4 Byte Length (72)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (TJMP) | 72bytes Unknown|
|---|---|---|---|---|
|48000000| 6E797361| 18BC2311 01000000 |706D6A74 | C0904402 01000000| 72 unknown bytes|


##### 3.3.5 Asyn FEED - CMSampleBuffer with h264 Video Data
For video data ASYN FEED packets, the device will use the ClockRef we sent as a reply to the CVRP Sync request. 
###### Packet Format Description

| 4 Byte Length (varies, 91607 in this example)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (FEED) | 4 bytes length of CMSampleBuf(varies, here 91587)| CMSampleBuf Magic (sbuf) | CMSampleBuf bytes|
|---|---|---|---|---|---|---|
|D7650100| 6E797361| 18BC2311 01000000 |64656566 | C3650100 | 66756273 | ... |



##### 3.3.6 Asyn EAT! - CMSampleBuffer with Audio Data
Just like FEED only with different Magic (0x21746165) and a CMSampleBuf containing audio. 

##### 3.3.7 Asyn NEED - Tell the device to send more
For telling the device to keep sending video data ASYN FEED packets, we need to send NEED packets with the ClockRef the device gave us in the  SYNC CVRP packet.
NEED Packets are constant over the whole session, so you can just init them once you received the correct clockRef and then just keep sending the same bytes over and over.
I think sending NEED packets is something you can do based on a timer (every 5 seconds f.ex.)
For easier implementation I just send one whenever I received a FEED.  
###### Packet Format Description

| 4 Byte Length (20)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID  |  4 bytes magic (NEED) |
|---|---|---|---|
|14000000| 6E797361| A08D5313 01000000 |6465656E | 

##### 3.3.8 Asyn HPD0 - Tell the device to stop video streaming
Send this to stop the device from streaming

###### Packet Format Description

| 4 Byte Length (20)   |4 Byte Magic (ASYN)   | 8 Byte empty clock CFTypeID  |  4 bytes magic (HPD0) |
|---|---|---|---|
|14000000| 6E797361| 01000000 00000000 |30617068 | 

##### 3.3.9 Asyn HPA0 - Tell the device to stop audio streaming

Send this to stop the device from streaming

###### Packet Format Description

| 4 Byte Length (20)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID of audio clock on device  |  4 bytes magic (HPA0) |
|---|---|---|---|
|14000000| 6E797361| 10FCC502 01000000 |30617068 | 

##### 3.3.10 Asyn RELS - Tell us about a released Clock on the device

###### Packet Format Description

| 4 Byte Length (20)   |4 Byte Magic (ASYN)   | 8 Byte clock CFTypeID of audio clock on device  |  4 bytes magic (RELS) |
|---|---|---|---|
|14000000| 6E797361| 008A6035 BA7F0000 | 736C6572 | 

## 4 Serializing/Deserializing Objects
### 4.0 General Description
This chapter explains how to serialize and deserialize all the necessary payload objects that you will find in the various SYNC and ASYN packets. 
### 4.1 Dictionaries
#### 4.1.0 General Description
Dictionaries are used throughout the protocol so they are pretty important to get right :-D
They are pretty easy to implement however note that there are two distinct types. Some dictionaries use only strings as keys and others use only ints or index numbers as keys.
Sometimes you will see dictionaries with other magic markers or single key value entries, but they always work the same way.
#### 4.1.1 Dictionaries with String Keys
##### 4.1.1.1 General Dictionary Structure

Dictionaries always start with a length int, dict magic which is then followed by a number of key value pairs each starting with a length field and keyv magic.
Every entry has a key starting with a 4byte int keylength and then followed by a strk magic int. Finally a string with the actual key.
Values work the same way as they start with a length, then a magic and the actual value. 
This example of a string key dictionary containing one boolean value nicely illustrates how dictionaries work. 

| 4 Byte Length (40)   |4 Byte Magic (DICT)   | 4byte length of first key value pair |  4 bytes magic of first key value pair (KEYV) | 4 byte length of first key(15)| stringkey magic (strk)| key string (Valeria) | 4 byte length of value(9)| 4byte value type magic (bulv==boolean) | value (0x1 == true) |
|---|---|---|---|---|---|---|---|---|---|
|28000000| 74636964| 20000000 |7679656B | 0F000000|6B727473|56616C65 726961|09000000|766C7562| 01|

Here are the value types I know about:

| magic little endian| magic big endian | description | value example |
|---|---|---|---|
|vlub|bulv|Boolean|0x1 or 0x0|
|vrts|strv|String|BlaBla|
|vtad|datv|Byte Array|0x010203|
|vbmn|nmbv|NSNumber|[4.4 NSNumber](#44-nsnumber)|
|tcid|dict|String or Index Key dict| see above |
|csdf|fdsc|CMFormatDescription|[4.5 CMFormatDescription](#45-cmformatdescription)|


#### 4.1.2 Dictionaries with 4-Byte Integer Index Keys
They work the same way as String Key Dictionaries with the only difference that all keys are 4 byte integers and they have 0x6B786469 (idxk) as a magic marker.
### 4.2 CMTime
This is exactly like in the CMTime.h https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.8.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMTime.h
And there is plenty of documentation so check it out there :-D

#### Example
|CMTimeValue  |CMTimeScale (Nanoseconds is the default)   | CMTimeFlags |  CMTimeEpoch |
|---|---|---|---|
|CB448EA1 CF10CC15| 00CA9A3B| 01000000 |00000000 00000000 | 


### 4.3 CMSampleBuffer
### 4.4 NSNumber
A very simple represenation of what probably is a NSNumber.
I have seen three different types:
- Type 3, 32 bit Integer
- Type 4, 64 bit Integer
- Type 6, 64 bit Float

#### Example
|4 byte Magic (nmbv)  |4 Byte int type   | Number, either 4 or 8 bytes depending on type |
|---|---|---|
|76626D6E | 03000000| 01000000 |



### 4.5 CMFormatDescription
Check out https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.9.sdk/System/Library/Frameworks/CoreMedia.framework/Versions/A/Headers/CMFormatDescription.h


## 5. Clocks and CMSync
I think the references in ASYN and SYNC packets are for CMClocks. So for sending a CMTime request I just a monotonic (DO NOT USE WALLCLOCK TIME) clock
to send a time difference in nanoseconds (Scale == 1000000000). It seems to work fine :-D
