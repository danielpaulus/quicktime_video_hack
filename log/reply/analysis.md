# Analysis of SYNC packets
All SYNC packets require us to reply with a RPLY packet. 
It seems like this is used for synchronizing what would be a few CMClock's (implement CMSync.h protocol)
on MacOSX.

## packet structure
req:
4 byte length
4 byte sync marker
10(4-4-2) bytes 01000000 000000 0030
4 byte cwpa magic
8 byte correlation 
6 bytes payload?

resp:
4 byte length
4 byte reply magic
8 bytes correlation
2 bytes 0030 aus dem req
10 bytes payload? oder 4 bytes 0 und 6 bytes payload

## How they work
It seems like the cwpa and cvrp sync packets are supposed to tell us to create CMClocks and send
back a referenceID for those.
I could observe that all subsequent Sync and Asyn Packets contain the same reference we send as a reply. 
So they probably tell us which CMClock to use for synching :-D 

## Different Sync Packet Types

This is an example list of packets received from the device in the exact order they appear
in the hexdump and what i currently think they could mean

|sync type   |meaning   | reply  |   |   |
|---|---|---|---|---|
|cwpa   |create clock, maybe for audio   | contains a 6byte reference to the created clock  |   |   |
|afmt(lpcm)   | probably audio format info   |  dict with error code 0 |   |   |
|cvrp   | create clock (maybe for video, the id is contained in all feed asyn packets)   | contains a 6byte reference to the created clock  |   |   |
|clok   |   |   |   |   |
|time   |   |   |   |   |
|time   |   |   |   |   |
|3x skew   |   |   |   |   |

## Details

### 1. CWPA Packet and Response

#### Example Request

| 4 Byte Length (36)   |4 Byte Magic (SYNC)   | 8 Empty clock reference| 2 bytes stuff(seems like a ID for sth.)  | 4 byte message type (CWPA)   | 8 byte correlation id  | 6 bytes identifier of the device clock |
|---|---|---|---|---|---|---|
|24000000 |636E7973 |01000000 00000000| 0030 | 61707763 |E03D5713 01000000| E074 5A130040 |

#### Example Response

Seems like the first two bytes of our clock identifier are always 0, and in the later packets appended so `0000 B00C E26CA67F` becomes  `B00C E26CA67F 0000`
| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |  Seems like the ID from the req. + two 0 bytes | 8 bytes identfier of our clock |
|---|---|---|---|---|---|
|1C000000 | 796C7072 |E03D5713 01000000 | 00300000 |0000 B00C E26CA67F |

### 2. AFMT Packet

#### Example Request

| 4 Byte Length (68)   |4 Byte Magic (SYNC)   | 8  clock reference| 4 byte magic (AFMT)| 8 byte correlation id| some weird data| 4 byte magic (LPCM) |  28 bytes what i think is pcm data|
|---|---|---|---|---|---|---|---|
|44000000| 636E7973| B00CE26C A67F0000| 746D6661 | 809D2213 01000000| 00000000 0070E740 |6D63706C| 4C000000 04000000 01000000 04000000 02000000 10000000 00000000|

#### Example Response
The response is basically a dictionary containing an error code, 0 if everything is ok :-D

| 4 Byte Length (62)   |4 Byte Magic (RPLY)   | 8  correlation id| 4 byte 0| 4 byte dict length(42)| 4 byte magic (DICT)| dict bytes |
|---|---|---|---|---|---|---|
|3E000000 |796C7072| 809D2213 01000000 |00000000| 2A000000| 74636964| 22000000 7679656B 0D000000 6B727473 4572726F 720D0000 0076626D 6E030000 0000|

### 3. CVRP Packet

#### Example Request

Contains a Dict with a FormatDescription and timing information
|4 Byte Length (649)|4 Byte Magic (SYNC)|8 byte empty(?) clock reference|4 byte magic(CVRP)|8 byte correlation id|reference id of clock on device (needs to be in NEED packets we send)|4 byte length of dictionary (613)|4 byte magic (DICT)| Dict bytes|
|---|---|---|---|---|---|---|---|---|
|89020000 |636E7973| 01000000 00000000 |70727663| D0595613 01000000 |A08D5313 01000000 |65020000| 74636964|   0x.....|

#### Example Response

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |  4 bytes (seem to be always 0) | 8 bytes identfier of our clock(will be in all feed async packets) |
|---|---|---|---|---|---|
|1C000000 | 796C7072 |D0595613 01000000 | 00000000 |5002D16C A67F0000 |


### 4. CLOK Packet
I am not quite sure what this is for, it seems like i am supposed to create a clock to then use it when sending two responses to time requests. 
Could be wrong though. 

#### Example Request

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock reference  |  4 bytes magic (CLOK) | 8 bytes correlation id |
|---|---|---|---|---|---|
|1C000000| 636E7973| 5002D16C A67F0000| 6B6F6C63 | 70495813 01000000 |

#### Example Response

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 correlation id  |  4 bytes (seem to be always 0) | 8 bytes identfier of our clock(for the next two time packets) |
|---|---|---|---|---|---|
|1C000000| 796C7072| 70495813 01000000| 00000000 | 8079C17C A67F0000|

### 5. TIME Packet

#### Example Request

| 4 Byte Length (28)   |4 Byte Magic (SYNC)   | 8 Byte clock reference  |  4 bytes magic (TIME) | 8 bytes correlation id |
|---|---|---|---|---|---|
|1C000000| 636E7973| 8079C17C A67F0000 |656D6974 | 503D2213 01000000 |



#### Example Response

| 4 Byte Length (44)   |4 Byte Magic (RPLY)   | 8 Byte correllation id  |  4 bytes 0x0 | 24 bytes CMTime struct |
|---|---|---|---|---|---|
|2C000000 |796C7072 |503D2213 01000000| 00000000 | E1E142C4 62BA0000 00CA9A3B 01000000 00000000 00000000|

### 6. SKEW Packet
