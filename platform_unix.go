// +build !windows

package main

import (
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
)

func openAgentSocket() (io.ReadWriteCloser, error) {
	path := os.Getenv(authSockEnv)
	if len(path) == 0 {
		return nil, errors.New(authSockEnv + " not defined")
	}
	return net.Dial("unix", path)
}

func setupCommandForPlatform(cmd *exec.Cmd) {
}
