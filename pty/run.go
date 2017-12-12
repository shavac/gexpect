package pty

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func Start(c *exec.Cmd) (term *Terminal, err error) {
	if term, err = NewTerminal(); err != nil {
		return nil, err
	}
	return term, term.Start(c)
}

func (t *Terminal) Start(c *exec.Cmd) (err error) {
	if t == nil {
		return errors.New("terminal not assigned.")
	}

	var stdout bytes.Buffer
	fmt.Println(t.Tty.Name())
	c.Stdout = bufio.NewWriter(&stdout)
	c.Stdin = t.Tty
	c.Stderr = bufio.NewWriter(&stdout)

	go func() {
		for {
			time.Sleep(10)
			by, _ := stdout.ReadBytes(20)
			t.Tty.Write(by)
			if t.Log != nil {
				t.Log.Write(by)
			}
		}
	}()

	c.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err = c.Start(); err != nil {
		fmt.Println("error is ", err)
		t.Pty.Close()
		return
	}
	return
}
