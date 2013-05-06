package gexpect

//import "fmt"

type Plan struct {
	VarList map[string]string
	Flow    Flow
}

func (p *Plan) BindVar(varname, value string) error {
	p.VarList[varname] = value
	return nil
}

func (p *Plan) VarApply() error {
	fn := func(step *Step) error {
		switch s:= (*step).(type) {
		case *VarSendStep:
			if v, ok := p.VarList[s.VarName]; ok {
				*step = SendStep{v}
			} else {
				return ValueNotFoundError{s.VarName}
			}
		case *VarSendLineStep:
			if v, ok := p.VarList[s.VarName]; ok {
				*step = SendStep{v + "\n"}
			} else {
				return ValueNotFoundError{s.VarName}
			}
		}
		return nil
	}
	return Walk(p.Flow, fn)
}

func NewPlan() *Plan{
	return &Plan{ make(map[string]string), Flow{}}
}
