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
in the hexdump

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
### 3. CVRP Packet

#### Example Request

Contains a Dict with a FormatDescription and timing information
|4 Byte Length (649)|4 Byte Magic (SYNC)|8 byte empty(?) clock reference|4 byte magic(CVRP)|8 byte correlation id|reference id of clock on device (needs to be in NEED packets we send)|4 byte length of dictionary (613)|4 byte magic (DICT)| Dict bytes|
|---|---|---|---|---|---|---|---|---|
|89020000 |636E7973| 01000000 00000000 |70727663| D0595613 01000000 |A08D5313 01000000 |65020000| 74636964|   0x.....|

#### Example Response

| 4 Byte Length (28)   |4 Byte Magic (RPLY)   | 8 Byte correlation id  |  4 bytes of stuff | 8 bytes identfier of our clock(will be in all feed async packets) |
|---|---|---|---|---|---|
|1C000000 | 796C7072 |D0595613 01000000 | 00000000 |5002D16C A67F0000 |