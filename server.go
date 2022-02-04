package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	// bufio reader to read data from the client
	netData := bufio.NewReader(conn)

	// buffers for different ip address
	ipv4Buffer := make([]byte, 4)
	ipv6Buffer := make([]byte, 16)

	// buffers for different kinds of ip addresses
	log.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		messageType, err := netData.ReadByte()

		if err != nil {
			conn.Close()
		}

		var addr net.IP

		switch messageType {
		case 0x04:
			// ipv4 address, read 4 bytes
			n, err := io.ReadFull(netData, ipv4Buffer)

			if n != 4 || err != nil {
				log.Println(err)
				continue
			}

			addr = net.IPv4(ipv4Buffer[0], ipv4Buffer[1], ipv4Buffer[2], ipv4Buffer[3])
		case 0x06:
			// ipv6 address, read 16 bytes
			n, err := io.ReadFull(netData, ipv6Buffer)

			if n != 16 || err != nil {
				log.Println(err)
				continue
			}

			addr = ipv6Buffer
		default:
			continue
		}

		// do geo ip lookup
		city := getGeoIP(addr)

		if city == nil {
			continue
		}

		// marshal city to json
		response, err := json.Marshal(city)

		if err != nil {
			log.Println(err)
			continue
		}

		// write response to client
		_, err = conn.Write(append(response, '\n'))

		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
	}
}
