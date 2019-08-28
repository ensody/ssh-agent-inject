package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ensody/ssh-agent-inject/common"
)

var (
	verbose = flag.Bool("v", false, "verbose output on stderr")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "Error: No positional arguments allowed.")
		flag.Usage()
		os.Exit(2)
	}

	path := os.Getenv(common.AuthSockEnv)
	if len(path) == 0 {
		log.Fatalln(common.AuthSockEnv + " not defined")
	}

	os.Remove(path)

	l, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalln("Listen error:", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		if *verbose {
			log.Println("Signal:", sig)
		}
		l.Close()
		os.Exit(0)
	}(sigc)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go forwardAgent(conn)
	}
}

var stdioLock = &sync.RWMutex{}

// forwardAgent forwards the given connection to the real ssh-agent via stdio.
func forwardAgent(conn net.Conn) {
	// This forwards each request & response pair in a loop between the connection and stdio.
	// The stdio part of the communication uses a lock to allow multiple clients to communicate in parallel.
	defer conn.Close()
	if *verbose {
		log.Println("Client connected")
		defer log.Println("Client disconnected")
	}
	for {
		if *verbose {
			log.Println("Reading request from connection")
		}
		packet, err := readAgentPacket(conn)
		if err != nil {
			if *verbose {
				log.Println("Error reading request from connection:", err)
			}
			return
		}
		packet, err = sendAgentRequest(packet)
		if err != nil {
			panic(err)
		}
		if *verbose {
			log.Println("Writing response to connection")
		}
		if _, err := conn.Write(packet); err != nil {
			log.Println("Error writing response to connection:", err)
			return
		}
		if *verbose {
			log.Println("Finished request-response sequence")
		}
	}
}

// sendAgentRequest sends an ssh-agent request and returns the respective ssh-agent response.
func sendAgentRequest(packet []byte) ([]byte, error) {
	if *verbose {
		log.Println("Acquiring stdio lock")
	}
	stdioLock.Lock()
	defer func() {
		if *verbose {
			log.Println("Releasing stdio lock")
		}
		stdioLock.Unlock()
	}()
	if *verbose {
		log.Println("Writing request to stdout")
	}
	if _, err := os.Stdout.Write(packet); err != nil {
		if *verbose {
			log.Println("Error writing request to stdout:", err)
		}
		return nil, err
	}
	if *verbose {
		log.Println("Reading response from stdin")
	}
	packet, err := readAgentPacket(os.Stdin)
	if err != nil {
		if *verbose {
			log.Println("Error reading response from stdin:", err)
		}
		return nil, err
	}
	return packet, nil
}

const maxPacketSize = 16 << 20

// readAgentPacket reads a whole ssh-agent packet from the given io.Reader.
func readAgentPacket(r io.Reader) ([]byte, error) {
	var rawLength [4]byte
	if _, err := io.ReadFull(r, rawLength[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(rawLength[:])
	if length == 0 {
		return nil, fmt.Errorf("Packet size is 0")
	}
	if length > maxPacketSize {
		return nil, fmt.Errorf("Packet size of %d is too large", length)
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	packet := append(rawLength[:], data...)
	return packet, nil
}
