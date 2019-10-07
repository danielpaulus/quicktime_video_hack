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

markers:
cwpa  createWith ?  reply contains a 6byte reference to the created clock :-D
afmt (lpcm) probably audio format info, reply is  a dict with error code 0
cvrp seems like it creates another clok, maybe for video? you can find the identifier in all feed asyncs
clok
time
time
3x skew