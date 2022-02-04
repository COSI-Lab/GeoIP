package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

const MAX_MIND_URL string = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&suffix=tar.gz&license_key="

var MAXMIND_LICENSE_KEY string

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
)

func main() {
	godotenv.Load()

	MAXMIND_LICENSE_KEY = os.Getenv("MAXMIND_LICENSE_KEY")

	if MAXMIND_LICENSE_KEY == "" {
		log.Fatal("MAXMIND_LICENSE_KEY not set")
	}

	listener, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)

	if err != nil {
		log.Fatalf("Could not create TCP listener %s\n", err)
	}

	defer listener.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	go handleDatabases()

	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}
