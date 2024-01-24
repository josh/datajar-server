//go:build linux

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/coreos/go-systemd/v22/activation"
)

var baseURL string

func main() {
	value, ok := os.LookupEnv("DATAJAR_URL")
	if !ok {
		panic("DATAJAR_URL not set")
	}
	baseURL = value

	listeners, err := activation.Listeners()
	if err != nil {
		panic(err)
	}

	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
	}
	ln := listeners[0]

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	unixAddr, ok := conn.RemoteAddr().(*net.UnixAddr)
	if !ok {
		log.Printf("Failed to get peer name: %s", unixAddr.Name)
		return
	}

	unitName, credID, err := parsePeerName(unixAddr.Name)
	if err != nil {
		log.Printf("Failed to parse peer name: %s", unixAddr.Name)
		return
	}
	log.Printf("%s requesting '%s' credential", unitName, credID)

	url := fmt.Sprintf("%s/%s", baseURL, credID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Failed to get credential: %v", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(conn, resp.Body)
	if err != nil {
		log.Printf("Failed to write credential: %v", err)
		return
	}
}

func parsePeerName(s string) (string, string, error) {
	matches := regexp.MustCompile("^@.*/unit/(.*)/(.*)$").FindStringSubmatch(s)
	if matches == nil {
		return "", "", fmt.Errorf("Failed to parse peer name: %s", s)
	}
	return matches[1], matches[2], nil
}
