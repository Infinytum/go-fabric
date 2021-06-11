# go-fabric
Fabric allows the transportation of TCP and UDP connections over any arbitrary way.


## Packet Header - 64 Bytes
4 Bytes - Packet Length
16 Bytes - Original Source IP
4 Bytes - Original Source Port
36 Bytes - Reserved

## Packet Payload - 65507 Bytes at most
Well... data