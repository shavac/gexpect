package gexpect

import (
	"regexp"
	"time"
	"fmt"
)

type Step interface {
	Do(*SubProcess) error
}

type ExpectStep struct {
	expects     []*regexp.Regexp
	timeout     int
	delay       time.Duration
	matchOK     bool
	WhenMatched Flow
	WhenTimeout Flow
}

func (es ExpectStep) Do(sp *SubProcess) error {
	es.matchOK, err = sp.Expect(es.timeout, es.expects...)
	return err
}

type SendStep struct {
	S string
}

func (f *Flow) SendLine(s string) *SendStep {
	return f.Send(s + "\n")
}

func (f *Flow) Send(s string) *SendStep {
	ss := &SendStep{s}
	*f = append(*f, ss)
	return ss
}

func (ss SendStep) Do(sp *SubProcess) error {
	return sp.Write(ss.S)
}

type VarSendStep struct {
	VarName string
}

func (s VarSendStep) Do(sp *SubProcess) error {
	return ValueNotBindError{s.VarName}
}

type VarSendLineStep struct {
	VarName string
}

func (s VarSendLineStep) Do(sp *SubProcess) error {
	return ValueNotBindError{s.VarName}
}

type PromptStep struct {
	Message string
}

func (f *Flow) Prompt(msg string) error {
	p := &PromptStep{msg}
	*f = append(*f, p)
	return nil
}

func (ps PromptStep) Do(sp *SubProcess) error {
	fmt.Println(ps.Message)
	return nil
}

type TerminateStep struct {
	Message string
}

func (f *Flow) Terminate(msg string) error {
	t := &TerminateStep{msg}
	*f = append(*f, t)
	return nil
}

func (ts TerminateStep) Do(sp *SubProcess) error {
	if err := sp.Close(); err != nil {
		return err
	}
	return TerminatedError{"terminate according plan with message: " + ts.Message}
}

type InteractStep struct {
}

func (f *Flow) Interact() error {
	step := &InteractStep{}
	*f = append(*f, step)
	return nil
}
func (is InteractStep) Do(sp *SubProcess) error {
	time.Sleep(5000)
	return nil
}