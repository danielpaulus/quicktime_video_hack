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
|cwpa   |createWith (packetAUdio???) ?   | contains a 6byte reference to the created clock :-D  |   |   |
|afmt(lpcm)   | probably audio format info   |  dict with error code 0 |   |   |
|cvrp   |   |   |   |   |
|clok   |   |   |   |   |
|time   |   |   |   |   |
|time   |   |   |   |   |
|3x skew   |   |   |   |   |

