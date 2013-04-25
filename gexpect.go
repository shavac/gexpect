package gexpect

import (
	"github.com/shavac/gexpect/pty"
	"os"
)

var (
	err error
)

type Child struct {
	pty, tty *os.File
}

func (c *Child) Spawn(commandline string) (err error) {
	c.pty, c.tty, err = pty.Open()
	return
}

