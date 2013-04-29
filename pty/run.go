package pty

import (
	"os"
	"os/exec"
	"syscall"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func Start(c *exec.Cmd) (pty, tty *os.File, err error) {
	pty, tty, err = Open()
	if err != nil {
		return nil, nil, err
	}
	defer tty.Close()
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	err = c.Start()
	if err != nil {
		pty.Close()
		return nil, nil, err
	}
	return pty, tty, err
}

func (t *Terminal)Start(c *exec.Cmd) (err error) {
	if (t == nil) {
		if t, err = NewTerminal(); err != nil {
		return err
		}
	}
	defer t.Tty.Close()
	c.Stdout = t.Tty
    c.Stdin = t.Tty
    c.Stderr = t.Tty
    err = c.Start()
    if err != nil {
        t.Pty.Close()
        return
    }
	return
}

func (t *Terminal)Close() (err error) {
	err = t.Tty.Close()
	err = t.Pty.Close()
	return
}