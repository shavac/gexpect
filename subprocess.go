package gexpect

import (
	"bufio"
	"github.com/shavac/gexpect/pty"
	"io"
	"os"
	"os/signal"
	"os/exec"
	"regexp"
	"time"
	"syscall"
)

var (
	err error
)

type SubProcess struct {
	term            *pty.Terminal
	cmd             *exec.Cmd
	DelayBeforeSend time.Duration
	CheckInterval   time.Duration
	Before          []byte
	After           []byte
	Match           []byte
	echo            bool
}
func init() {
}

func (sp *SubProcess) Start() (err error) {
	return sp.term.Start(sp.cmd)
}

func (sp *SubProcess) Close() (err error) {
	sp.Terminate()
	return sp.term.Close()
}

func (sp *SubProcess) WaitTimeout(d time.Duration) (err error) {
	if sp.echo {
		go func() {
			io.Copy(os.Stdout, sp)
		}()
	}
	timeout := make(chan bool, 1)
	execerr := make(chan error, 1)
	go func() {
		if d == 0 {
			return
		}
		time.Sleep(d)
		timeout <- true
	}()
	go func() {
		execerr <- sp.cmd.Wait()
	}()
	select {
	case <-timeout:
		return TIMEOUT
	case err := <-execerr:
		return err
	}
}

func (sp *SubProcess) Wait() error {
	return sp.WaitTimeout(0)
}

func (sp *SubProcess) Terminate() error {
	return sp.cmd.Process.Kill()
}

func (sp *SubProcess) Expect(timeout time.Duration, expreg ...*regexp.Regexp) (matchIndex int, err error) {
	buf := make([]byte, 2048)
	c := make(chan byte, 1)
	tmout := make(chan bool, 1)
	checkpoint := make(chan int, 1)
	rerr := make(chan error, 1)
	go func() {
		for {
			b := make([]byte, 1)
			if _, err := io.ReadAtLeast(sp, b, 1); err != nil {
				rerr <- err
			}
			c <- b[0]
			if sp.echo {
				os.Stdout.Write(b)
				os.Stdout.Sync()
			}
		}
	}()
	go func() {
		for i := 1; ; i++ {
			time.Sleep(time.Microsecond)
			checkpoint <- i
		}
	}()
	go func() {
		time.Sleep(timeout)
		tmout <- true
	}()
	for {
		select {
		case c1 := <-c:
			buf = append(buf, c1)
		case <-tmout:
			sp.Before = append(sp.Before, buf...)
			return -1, TIMEOUT
		case e := <-rerr:
			sp.Before = append(sp.Before, buf...)
			return -1, e
		case <-checkpoint:
			for idx, re := range expreg {
				if loc := re.FindIndex(buf); loc != nil {
					sp.Match = buf[loc[0]:loc[1]]
					sp.Before = append(sp.Before, buf[0:loc[0]]...)
					buf = make([]byte, 2048)
					return idx, nil
				}
			} // no match
		}
	}
	return -1, nil
}

func (sp *SubProcess) Read(b []byte) (n int, err error) {
	return sp.term.Read(b)
}

func (sp *SubProcess) Write(b []byte) (n int, err error) {
	time.Sleep(sp.DelayBeforeSend)
	return sp.term.Write(b)
}

func (sp *SubProcess) Writeln(b []byte) (n int, err error) {
	bn := append(b, []byte("\r\n")...)
	return sp.Write(bn)
}

func (sp *SubProcess) Send(response string) (err error) {
	_, err = sp.Write([]byte(response))
	return
}

func (sp *SubProcess) SendLine(response string) (err error) {
	return sp.Send(response + "\r\n")
}

func (sp *SubProcess) Interact() (err error) {
	return sp.InteractTimeout(0)
}

func (sp *SubProcess) InteractTimeout(d time.Duration) (err error) {
	sp.Write(sp.After)
	sp.After = []byte{}
	sp.term.SetRaw()
	defer sp.term.Restore()
	s := make(chan os.Signal, 1)
    signal.Notify(s, os.Interrupt, syscall.SIGWINCH)
    go func() {
		for sig := range s {
            switch sig {
			case os.Interrupt:
				sp.term.SendIntr()
			case syscall.SIGWINCH:
				sp.term.ResetWinSize()
            default:
                continue
            }
        }
    }()
	timeout := make(chan bool, 1)
	go func() {
		if d == 0 {
			return
		}
		time.Sleep(d)
		timeout <- true
	}()
	execerr := make(chan error, 1)
	go func() {
		execerr <- sp.cmd.Wait()
	}()
	in := make(chan byte, 1)
	stdin := bufio.NewReader(os.Stdin)
	go func() error {
		var b byte
		for {
			if b, err = stdin.ReadByte(); err != nil {
				if err == io.EOF {
					sp.term.SendEOF()
					continue
				} else {
					return err
				}
			}
			in <- b
		}
	}()
	go func() {
		io.Copy(os.Stdout, sp)
		return
	}()
	for {
		select {
		case <-timeout:
			return TIMEOUT
		case err := <-execerr:
			return err
		case b := <-in:
			_, err = sp.Write([]byte{b})
		}
	}
	return
}

func (sp *SubProcess) Echo() {
	sp.echo = true
}

func (sp *SubProcess) NoEcho() {
	sp.echo = false
}

func NewSubProcess(name string, arg ...string) (sp *SubProcess, err error) {
	sp = new(SubProcess)
	sp.term, err = pty.NewTerminal()
	sp.cmd = exec.Command(name, arg...)
	sp.DelayBeforeSend = 50 * time.Microsecond
	sp.CheckInterval = time.Microsecond
	return
}
