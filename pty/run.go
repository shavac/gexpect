package pty

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

var (
	stdout bytes.Buffer
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

	//defer stdout.Reset()

	c.Stdout = bufio.NewWriter(&stdout)
	//c.Stdout = t.Tty
	//c.Stdin = t.Tty
	//c.Stderr = t.Tty
	c.Stdin = t.Tty
	c.Stderr = bufio.NewWriter(&stdout)

	c.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err = c.Start(); err != nil {
		fmt.Println("error is ", err)
		t.Pty.Close()
		return
	}

	//frd := bufio.NewReader(&stdout)
	//ch := make(chan bool, 1)
	//tee := io.TeeReader(t.Tty, t.Log)
	go func() {
		//w := make([]byte, 4096)

		for {
			by, _ := stdout.ReadBytes(10)
			if by == nil {
				continue
			}
			//w = w[:cap(w)]
			//n, err :=
			//stdout.WriteByte(w)
			//_, err := stdout.WriteTo(w)
			//if err != nil {
			//	fmt.Println(err)
			//}
			// io.TeeReader(t.Tty, t.Log)

			//nb := stdout.Len()
			//if nb > 50 {
			//	-nb = 50
			//}

			//w := stdout.Next(nb)
			//s, err := stdout.Read(w)
			//scanner := bufio.NewScanner(&stdout)
			//for scanner.Scan() {
			//	t.Tty.Write(scanner.Bytes())
			//	if t.Log != nil {
			//		t.Log.Write(scanner.Bytes())
			//	}
			//}
			//fmt.Println(s)
			//if s == 0 {
			//	continue
			//}
			//fmt.Println(s)
			//for scanner.
			//_, err :=
			//	io.ByteReader
			//	by, err := rd.ReadBytes('\n')
			//if r > 0 {
			//	fmt.Println(by)
			//}
			//by, err := stdout.ReadBytes('\n')
			//if err != nil {
			//	fmt.Println(err)
			//	continue
			//}
			t.Tty.Write(by)
			if t.Log != nil {
				t.Log.Write(by)
			}

		}
		//io.Copy(t.Log, t.Tty)
	}()

	return
}
