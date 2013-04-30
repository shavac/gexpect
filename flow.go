package gexpect

import (
	"regexp"
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

