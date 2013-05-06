package gexpect

import (
	"regexp"
	//"fmt"
)

type Flow []Step

func (f *Flow) VarSend(varname string) *VarSendStep {
	sa := &VarSendStep{varname}
	*f = append(*f, sa)
	return sa
}

func (f *Flow) VarSendLine(argname string) *VarSendLineStep {
    vsl := &VarSendLineStep{argname}
    *f = append(*f, vsl)
    return vsl
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


func Walk(flow Flow, fn func(step *Step) error) error {
	for i, _ := range flow {
		if err := fn(&flow[i]); err != nil {
			return err
		}
		switch s := flow[i].(type) {
		case *ExpectStep:
			if err := Walk(s.WhenMatched, fn); err != nil {
				return err
			}
			if err := Walk(s.WhenTimeout, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
