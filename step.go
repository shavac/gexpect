package gexpect

import (
	"regexp"
	"time"
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
}

type TerminateStep struct{}

func (f *Flow) Terminate() error {
	t := new(TerminateStep)
	*f = append(*f, t)
	return nil
}

func (ts TerminateStep) Do(sp *SubProcess) error {
	return sp.Close()
}