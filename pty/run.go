package pty

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

	//x := io.MultiWriter(t.Tty, &stdout)
	//c.Stdout = bufio.NewWriter(&stdout)
	//c.Stdout = bufio.NewWriterSize(&stdout, 10485760)
	fmt.Println(t.Log.Name())

	c.Stdout = io.MultiWriter(t.Tty, t.Log)

	//c.Stdout = t.Tty
	c.Stdin = t.Tty
	//c.Stderr = bufio.NewWriter(&stdout)
	//c.Stderr = bufio.NewWriterSize(&stdout, 10485760)
	c.Stderr = io.MultiWriter(t.Tty, t.Log)
	//c.Stderr = t.Tty

	c.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	if err = c.Start(); err != nil {
		fmt.Println("error is ", err)
		t.Pty.Close()
		return
	}

	//go func() {
	//	for {
	//stdr := bufio.NewReader(&stdout)
	//w := make([]byte, 4096)
	//	time.Sleep(100)
	//			by, _ := stdout.ReadBytes(('\x0a' |
	//				'\x23' |
	//				'\x24' |
	//				'\x25' |
	//				'\x22' |
	//				'\x27' |
	//				'\x7b' |
	//				'\x7d' |
	//				'\x09' |
	//				'\x0b' |
	//				'\x20'))
	//		by, _ := stdout.ReadBytes('\n' | '\'' | '"' | ' ' | '\\')
	//by := stdout.Next(256)

	//by, _ := stdout.ReadBytes(' ')

	//reader := bufio.NewReader(&stdout)
	//by, _ := ioutil.ReadAll(reader)

	//e_, err := io.ReadFull(stdr, w)

	//if err != io.EOF {
	//	fmt.Fprintln(os.Stderr, err)
	//}
	//by := make([]byte, 4096)
	//n, err := stdout.Read(by)
	//if n == 0 {
	//	continue
	//}
	//if err == io.EOF {
	//	continue
	//}
	//by = []byte{}
	//for i := 1; i <= 4096; i++ {
	//	c, err := stdout.ReadByte()
	//	if err == io.EOF {
	//		break
	//	}
	//	by = append(by, c)
	//}

	//	by, err := stdout.ReadBytes((' ' | '\n'))
	//	if err != io.EOF && err != nil {
	//		fmt.Println(err)
	//	}

	//t.Tty.Write(by)
	//		if t.Log != nil {
	//			t.Log.Write(by)
	//		}
	//by = []byte{}
	//	}
	//}()

	return
}
