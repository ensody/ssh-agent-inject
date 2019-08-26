package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	verbose = flag.Bool("v", false, "verbose output on stderr")
)

const authSockEnv = "SSH_AUTH_SOCK"

func main() {
	flag.Parse()

	path := os.Getenv(authSockEnv)
	if len(path) == 0 {
		log.Fatalln(authSockEnv + " not defined")
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

	// Dynamically pass stdin (i.e. the responses from the host's ssh-agent) to the current
	// connected client
	conn := &dynamicWriter{}
	go func() {
		io.Copy(conn, os.Stdin)
		l.Close()
		os.Exit(0)
	}()

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		if *verbose {
			log.Println("Client connected")
		}
		conn.Lock()
		conn.conn = c
		conn.Unlock()

		io.Copy(os.Stdout, c)
		if *verbose {
			log.Println("Client disconnected")
		}

		conn.Lock()
		c.Close()
		conn.conn = nil
		conn.Unlock()
	}
}

type dynamicWriter struct {
	sync.RWMutex
	conn io.Writer
}

func (d *dynamicWriter) Write(b []byte) (int, error) {
	for {
		d.Lock()
		defer d.Unlock()
		if d.conn == nil {
			if *verbose {
				log.Println("Discarding write")
			}
			return len(b), nil
		}
		n, err := d.conn.Write(b)
		if err == nil {
			if *verbose {
				log.Printf("Wrote %d bytes\n", n)
			}
			return n, err
		}
		if *verbose {
			log.Println("Error during write:", err)
		}
	}
}
