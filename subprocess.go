package gexpect

import (
	"github.com/shavac/gexpect/pty"
	"os/exec"
	"regexp"
)

var (
	err error
)

type SubProcess struct {
	term *pty.Terminal
	cmd *exec.Cmd
}

func (sp *SubProcess) Start() (err error) {
	return nil
}

func (sp *SubProcess) Close() (err error) {
	return nil
}

func (sp *SubProcess) Expect(timeout int, expreg ...*regexp.Regexp) (matchOK bool, err error) {
	return true, nil
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
	sp.term, err = pty.NewTerminal()
	sp.cmd = exec.Command(name, arg...)
	return
}

