# GeoIP

This was a project I worked on to create a networked service to preform GeoIP lookups for Mirror and other services. By working on this project I learned that maxmind is actually very generious and there isn't a large need for a service like this. 

## Protocol

Message types

| 1st Byte        | Description  |
| --------------- | ------------ |
| `0x04`          | ipv4 address |
| `0x06`          | ipv6 address |
| Everything else | Unused       |
