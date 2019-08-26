package main

import (
	"syscall"
	"io"
	"os"
	"os/exec"

	"github.com/Microsoft/go-winio"
)

const sshAgentPipe = "//./pipe/openssh-ssh-agent"

func openAgentSocket() (io.ReadWriteCloser, error) {
	conn, err := winio.DialPipe(sshAgentPipe, nil)
	if err != nil {
		err = &os.PathError{Path: sshAgentPipe, Op: "open", Err: err}
	}
	return conn, err
}

func setupCommandForPlatform(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
