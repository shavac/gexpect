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
	term *pty.Terminal
	cmd *exec.Cmd
	DelayBeforeSend time.Duration
}

func (sp *SubProcess) Start() (err error) {
	return sp.term.Start(sp.cmd)
}

func (sp *SubProcess) Close() (err error) {
	return sp.term.Close()
}

func (sp *SubProcess) Expect(timeout time.Duration, expreg ...*regexp.Regexp) (matchOK bool, err error) {
	/*var buf []byte
	timeout := make(chan bool, 1)
	checkpoint := make(chan int, 1)
	 */
	return true, nil
}

func (sp *SubProcess) Read(b []byte) (n int, err error) {
	return sp.term.Read(b)
}

func (sp *SubProcess) Write(response string) (err error) {
	_, err = sp.term.Write([]byte(response))
	return
}

func (sp *SubProcess) Writeln(response string) (err error) {
	return sp.Write(response+"\n")
}

func (sp *SubProcess) Interact() (err error) {
	return nil
}

func NewSubProcess(name string, arg ...string) (sp *SubProcess, err error) {
	sp := new(SubProcess)
	sp.term, err = pty.NewTerminal()
	sp.cmd = exec.Command(name, arg...)
	sp.DelayBeforeSend = 50 * time.Microsecond
	return
}

