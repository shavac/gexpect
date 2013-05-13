package gexpect

import (
	"github.com/shavac/gexpect/pty"
	"os/exec"
	"regexp"
	"time"
)

var (
	err error
)

type SubProcess struct {
	term            *pty.Terminal
	cmd             *exec.Cmd
	DelayBeforeSend time.Duration
	CheckInterval   time.Duration
	before          []byte
	after           []byte
	match           []byte
}

func (sp *SubProcess) Start() (err error) {
	return sp.term.Start(sp.cmd)
}

func (sp *SubProcess) Close() (err error) {
	return sp.term.Close()
}

func (sp *SubProcess) Expect(timeout time.Duration, expreg ...*regexp.Regexp) (matchOK bool, err error) {
	var buf []byte
	tmout := make(chan bool, 1)
	checkpoint := make(chan int, 1)
	rerr := make(chan error, 1)
	go func() {
		for {
			_, err := sp.Read(buf)
			rerr <- err
		}
	}()
	go func() {
		for i := 1; ; i++ {
			time.Sleep(time.Microsecond)
			checkpoint <- i
		}
	}()
	go func() {
		time.Sleep(sp.CheckInterval)
		tmout <- true
	}()
	for {
		select {
		case <-tmout:
			sp.before = append(sp.before, buf...)
			return false, TIMEOUT
		case err := <-rerr:
			sp.before = append(sp.before, buf...)
			return false, err
		case <-checkpoint:
			for _, re := range expreg {
				if sp.match = re.Find(buf); sp.match != nil {
					ind := re.FindIndex(buf)
					sp.before = append(sp.before, buf[0:ind[0]]...)
					sp.after = append(sp.after, buf[ind[0]+ind[1]:]...)
					return true, nil
				}
			} // no match
			sp.before = append(sp.before, buf...) //fix this
			buf = []byte{}
		}
	}
	return false, nil
}

func (sp *SubProcess) Read(b []byte) (n int, err error) {
	return sp.term.Read(b)
}

func (sp *SubProcess) Write(response string) (err error) {
	time.Sleep(sp.DelayBeforeSend)
	_, err = sp.term.Write([]byte(response))
	return
}

func (sp *SubProcess) Writeln(response string) (err error) {
	return sp.Send(response + "\n")
}

func (sp *SubProcess) Send(response string) (err error) {
	return sp.Write(response)
}

func (sp *SubProcess) SendLine(response string) (err error) {
	return sp.Writeln(response)
}

func (sp *SubProcess) Interact() (err error) {
	return nil
}

func NewSubProcess(name string, arg ...string) (sp *SubProcess, err error) {
	sp = new(SubProcess)
	sp.term, err = pty.NewTerminal()
	sp.cmd = exec.Command(name, arg...)
	sp.DelayBeforeSend = 50 * time.Microsecond
	sp.CheckInterval = time.Microsecond
	return
}
