# Technical Documentation of the iOS ScreenSharing Feature
## 0. General Information
This document provides you with details about the screen sharing feature of QuickTime for iOS devices. 
The information contained in this document can be used to re-implement that feature in the programming language of choice and use
the feature on other operating systems than MAC OS X.
The repository also contains a reference implementation in Golang. 

-- Note: All the information in this document is reverse engineered by the author, therefore it could be wrong or not entirely accurate
as it involves a lot of assumptions. If you find mistakes or more accurate descriptions please add them :-) -- 


## 1. How to Enable it for a iOS Device on the USB Level
### 1.1 Foundations
Usually devices attached on the USB Port have a set of "configurations" that you can retrieve using the LibUsb wrapper you use.
Inside of these interfaces there are a set of Usb Endpoints you can use to communicate with your device. We are interested in the 
`bulk` endpoints of iOS devices as these are used for transferring data. 
### 1.2 USBMux Bulk Endpoints                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     
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
## 2. How to Initiate a Video Recording Session
## 3. Protocol Reference
### 3.1 Ping Packet
### 3.2 Sync Packets
#### 3.2.1 General Description
All SYNC packets require us to reply with a RPLY packet. 
It seems like this is mostly used for synchronizing CMClocks and exchanging 8byte CMClockRefs (implement CMSync.h protocol)
Usually you can see that SYNC packets have a 4byte SUB-TYPE followed by a 8byte correlationID. A reply always contains the correlationID so 
I assume this is how the device knows which reply belongs to which request.


##### 3.2.2. CWPA Packet and Response
##### General Description
This packet seems to be used for intitiating the audio stream. We get a clockRef from the device and respond with our own, newly created clockRef.
The clockref send by the device needs to go in the ASYN-1APH packet we send.
##### Request Format Description

| 4 Byte Length (36)   |4 Byte Magic (SYNC)   | 8 Empty clock reference| 4 byte message type (CWPA)   | 8 byte correlation id  | 8 bytes CFTypeID of the device clock |
|---|---|---|---|---|---|---|
|24000000 |636E7973 |01000000 00000000 | 61707763 |E03D5713 01000000| E0740000 5A130040 |

#### Example RPLY

Sends back our clockRef. The device will use the clockRef from here in the SYNC_AFMT message to tell us about the audio format. 
Also this will be used for all ASYN_EAT packets containing audio sample buffers. 

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |   4 Byte: 0  | 8 bytes CFTypeID of our clock |
|---|---|---|---|---|---|
|1C000000 | 796C7072 |E03D5713 01000000 | 00000000 |B00CE26C A67F0000|

##### 3.2.3. AFMT Packet
##### General Description
This packet contains information about the Audio Format(AMFT).  I assume the lpcm marker, which is followed by 7 integer values is somekind of information about Linear pulse-code modulation (LPCM)
The response is basically a dictionary containing an error code. Normally we send 0 to indicate everything is ok.
Note how the device references the Clock we gave it in the SYNC_CWPA_RPLY
#### Request Format Description

| 4 Byte Length (68)   |4 Byte Magic (SYNC)   | 8 bytes clock CFTypeID| 4 byte magic (AFMT)| 8 byte correlation id| some weird data| 4 byte magic (LPCM) |  28 bytes what i think is pcm data|
|---|---|---|---|---|---|---|---|
|44000000| 636E7973| B00CE26C A67F0000| 746D6661 | 809D2213 01000000| 00000000 0070E740 |6D63706C| 4C000000 04000000 01000000 04000000 02000000 10000000 00000000|

#### Reply(RPLY) Format Description
Contains the correlationID from the request as well as a simple Dictionary:  {"Error":NSNumberUint32(0)}

| 4 Byte Length (62)   |4 Byte Magic (RPLY)   | 8  correlation id| 4 byte 0| 4 byte dict length(42)| 4 byte magic (DICT)| dict bytes |
|---|---|---|---|---|---|---|
|3E000000 |796C7072| 809D2213 01000000 |00000000| 2A000000| 74636964| 22000000 7679656B 0D000000 6B727473 4572726F 720D0000 0076626D 6E030000 0000|

##### 3.2.4. CVRP Packet
##### General Description
The CVRP packet is used to create a local clock for synchronizing video frames. The ClockRef we get from the device, is the one we need to put in NEED packets.
Similar to this the device sends all ASYN_FEED sample buffer with a reference to our clock. After this, two SetProperty Asyns are usually received plus CLok, TBAS and SRAT.
#### Request Format Description

Contains a Dict with a FormatDescription and timing information. With the FormatDescription for the first time we get the h264 Picture and Sequence ParameterSet(PPS/SPS) already encoded in nice NALUs ready for streaming
over the network. They are hidden inside a dictionary inside the extension dictionary.
|4 Byte Length (649)|4 Byte Magic (SYNC)|8 byte empty(?) clock reference|4 byte magic(CVRP)|8 byte correlation id|CFTypeID of clock on device (needs to be in NEED packets we send)|4 byte length of dictionary (613)|4 byte magic (DICT)| Dict bytes|
|---|---|---|---|---|---|---|---|---|
|89020000 |636E7973| 01000000 00000000 |70727663| D0595613 01000000 |A08D5313 01000000 |65020000| 74636964|   0x.....|

#### Reply(RPLY) Format Description

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |  4 bytes (seem to be always 0) | 8 bytes CFTypeID of our clock(will be in all feed async packets) |
|---|---|---|---|---|
|1C000000 | 796C7072 |D0595613 01000000 | 00000000 |5002D16C A67F0000 |


##### 3.2.5. CLOK Packet
##### General Description
I am not quite sure what this is for, it seems like i am supposed to create a clock to then use it when sending two responses to time requests. 


#### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (CLOK) | 8 bytes correlation id |
|---|---|---|---|---|
|1C000000| 636E7973| 5002D16C A67F0000| 6B6F6C63 | 70495813 01000000 |

#### Reply(RPLY) Format Description

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 correlation id  |  4 bytes (seem to be always 0) | 8 bytes CFTypeID of our clock(for the next two time packets) |
|---|---|---|---|---|
|1C000000| 796C7072| 70495813 01000000| 00000000 | 8079C17C A67F0000|

##### 3.2.6. TIME Packet
##### General Description
#### Request Format Description

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock CFTypeID  |  4 bytes magic (TIME) | 8 bytes correlation id |
|---|---|---|---|---|
|1C000000| 636E7973| 8079C17C A67F0000 |656D6974 | 503D2213 01000000 |

#### Reply(RPLY) Format Description

| 4 Byte Length (44)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  |  4 bytes 0x0 | 24 bytes CMTime struct |
|---|---|---|---|---|
|2C000000 |796C7072 |503D2213 01000000| 00000000 | E1E142C4 62BA0000 00CA9A3B 01000000 00000000 00000000|

##### 3.2.7. SKEW Packet
##### General Description
### 3.2 Asyn Packets
## 4 Serializing/Deserializing Objects
### 4.0 Dictionaries
### 4.1 CMTime