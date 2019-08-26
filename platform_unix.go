// +build !windows

package main

import (
	"errors"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/ensody/ssh-agent-inject/common"
)

func openAgentSocket() (io.ReadWriteCloser, error) {
	path := os.Getenv(common.AuthSockEnv)
	if len(path) == 0 {
		return nil, errors.New(common.AuthSockEnv + " not defined")
	}
	return net.Dial("unix", path)
}

func setupCommandForPlatform(cmd *exec.Cmd) {
}
