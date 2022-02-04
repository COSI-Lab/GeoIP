package main

import (
	"bufio"
	"crypto/rand"
	"log"
	"net"
	"time"
)

var count int

// Read from TCP server and print to stdout
func read(reader *bufio.Reader) {
	for {
		count++

		// Read from TCP server
		_, err := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
			return
		}
	}
}

// Write random ipv4s to TCP server
func write(writer *bufio.Writer) {
	for {
		// Generate random ipv4
		ipv4 := make([]byte, 4)

		// random
		_, err := rand.Read(ipv4)

		if err != nil {
			log.Println(err)
			return
		}

		// Write to TCP server
		_, err = writer.Write([]byte{0x04})
		if err != nil {
			log.Println(err)
			return
		}

		_, err = writer.Write(ipv4)
		if err != nil {
			log.Println(err)
			return
		}

		// Flush the writer
		err = writer.Flush()

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	// Connect to TCP server
	conn, err := net.Dial("tcp", "localhost:3333")

	if err != nil {
		log.Fatal(err)
	}

	// split conn into reader and writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send and recieve data
	go read(reader)
	go write(writer)

	// Sleep for 1 minute
	time.Sleep(time.Minute)

	// close the connection
	conn.Close()

	log.Printf("%d requests, ave = %f requests per second", count, float64(count)/60)
}
