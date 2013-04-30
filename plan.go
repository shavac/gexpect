package gexpect

type Plan struct {
	VarList map[string]string
	Flow Flow
}

func (p *Plan) BindVar(varname, value string) error {
	p.VarList[varname]=value
	return nil
}

func (p *Plan) BindApply() error {
	var iter func(Flow) error
	iter = func(flow Flow) error {
		for i, step := range flow {
			switch s := step.(type) {
			case VarSendStep:
				if v, ok := p.VarList[s.VarName]; ok {
					flow[i] = SendStep{v}
				} else {
					return ValueNotFoundError{s.VarName}
				}
			case VarSendLineStep:
				if v, ok := p.VarList[s.VarName]; ok {
                    flow[i] = SendStep{v+"\n"}
                } else {
                    return ValueNotFoundError{s.VarName}
                }
			case ExpectStep:
				if err:= iter(s.WhenMatched); err != nil {
					return err
				}
				if err:= iter(s.WhenTimeout); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return iter(p.Flow)
}
