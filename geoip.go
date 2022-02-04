package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader
var db_lock sync.RWMutex

// handleDatabases checks for new databases and downloads them
func handleDatabases() {
	// Check for a new database every day
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		log.Println("Checking for new database")

		if !checkForNewDatabase() {
			log.Println("No new database found")
		} else {
			log.Println("New database found. Downloading...")
			err := downloadNewDatabase()

			if err != nil {
				log.Println("Error downloading new database:", err)
			} else {
				log.Println("Download complete. Opening...")
				err = openNewDatabase("GeoLite2-City.mmdb")

				if err != nil {
					log.Println("Error opening new database:", err)
				}
			}
		}
	}
}

// openNewDatabase goes to a file on disk and opens the database
// the old database is freed.
func openNewDatabase(filename string) error {
	newdb, err := geoip2.Open(filename)

	if err != nil {
		return err
	}

	db_lock.Lock()
	defer db_lock.Unlock()

	if db != nil {
		db.Close()
	}

	db = newdb

	return nil
}

func checkForNewDatabase() bool {
	// Make a HEAD request to MaxMind
	url, err := url.Parse(MAX_MIND_URL + MAXMIND_LICENSE_KEY)

	if err != nil {
		return false
	}

	req := http.Request{Method: http.MethodHead, URL: url}
	resp, err := http.DefaultClient.Do(&req)

	if err != nil {
		log.Println("Error checking for new database:", err)
		return false
	}

	// Get last modified header
	lastModified := resp.Header.Get("Last-Modified")

	// Load the last modified date from the file
	lastModifiedFile, err := os.Open("last_modified")

	if err != nil {
		// create the file
		lastModifiedFile, err = os.Create("last_modified")

		if err != nil {
			log.Println("Error creating last_modified file:", err)
			return false
		}
	}

	defer lastModifiedFile.Close()

	// Read the last modified date from the file
	lastModifiedFileBytes, err := io.ReadAll(lastModifiedFile)

	if err != nil {
		log.Println("Error reading last_modified file:", err)
		return false
	}

	return lastModified != string(lastModifiedFileBytes)
}

func downloadNewDatabase() error {
	// Download the database from maxmind
	url, err := url.Parse(MAX_MIND_URL + MAXMIND_LICENSE_KEY)

	if err != nil {
		return err
	}

	req := http.Request{Method: http.MethodGet, URL: url}
	resp, err := http.DefaultClient.Do(&req)

	if err != nil {
		return err
	}

	// Extract the tarball
	gzr, err := gzip.NewReader(resp.Body)

	if err != nil {
		return err
	}

	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		if err != nil {
			log.Println("Error reading tarball:", err)
			return err
		}

		if strings.Split(header.Name, "/")[1] == "GeoLite2-City.mmdb" {
			log.Println("Extracting database...")

			// write the file to disk
			f, err := os.Create("GeoLite2-City.mmdb")

			if err != nil {
				return err
			}

			io.Copy(f, tr)
			f.Close()
			break
		}
	}

	// Write the last modified date to the file
	lastModifiedFile, err := os.Create("last_modified")

	if err != nil {
		return err
	}

	_, err = lastModifiedFile.Write([]byte(resp.Header.Get("Last-Modified")))

	if err != nil {
		return err
	}

	lastModifiedFile.Close()

	return nil
}

// getGeoIP returns the geoip information for the given ip address
func getGeoIP(ip net.IP) *geoip2.City {
	// Acquire a read only lock to the database
	db_lock.RLock()
	city, error := db.City(ip)
	db_lock.RUnlock()

	if error != nil {
		return nil
	}

	return city
}
