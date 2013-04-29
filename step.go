package gexpect

import (
	"errors"
	"regexp"
)

type Step interface {
	Do(*SubProcess) error
}

type Flow []Step

type ExpectStep struct {
	expects     []*regexp.Regexp
	timeout     int
	matchOK     bool
	before      []byte
	after       []byte
	match       []byte
	WhenMatched Flow
	WhenTimeout Flow
}

func (es *ExpectStep) Do(sp *SubProcess) error {
	es.matchOK, err = sp.Expect(es.timeout, es.expects...)
	return err
}

func (f *Flow) Expect(timeout int, expects ...string) *ExpectStep {
	es := new(ExpectStep)
	es.timeout = timeout
	for _, e := range expects {
		if expregex, err := regexp.Compile(e); err != nil {
			return nil
		} else {
			es.expects = append(es.expects, expregex)
		}
	}
	*f = append(*f, es)
	return es
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

func (ss *SendStep) Do(sp *SubProcess) error {
	return sp.Write(ss.S)
}

type SendArgStep struct {
	Argname string
}

func (sas *SendArgStep) Do(sp *SubProcess) error {
	if v, ok := Args[sas.Argname]; ok {
		return sp.Write(v)
	} else {
		return errors.New("No such argument: " + sas.Argname)
	}
}

func (f *Flow) SendArg(argname string) *SendArgStep {
	sa := &SendArgStep{argname}
	*f = append(*f, sa)
	return sa
}

type SendArgLineStep struct {
    Argname string
}

func (sals *SendArgLineStep) Do(sp *SubProcess) error {
    if v, ok := Args[sals.Argname]; ok {
        return sp.Write(v+"\n")
    } else {
        return errors.New("No such argument: " + sals.Argname)
    }
}

func (f *Flow) SendArgLine(argname string) *SendArgLineStep {
    sal := &SendArgLineStep{argname}
    *f = append(*f, sal)
    return sal
}

type PromptStep struct {
}

type TerminateStep struct{}

func (f *Flow) Terminate() error {
	t := new(TerminateStep)
	*f = append(*f, t)
	return nil
}

func (ts *TerminateStep) Do(sp *SubProcess) error {
	return sp.Close()
}