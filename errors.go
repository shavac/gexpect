package gexpect

import (
)

type ValueNotBindError struct {
	VarName string
}

func (e ValueNotBindError) Error() string {
	return "Value not bind: " + e.VarName
}

type ValueNotFoundError struct {
	VarName string
}

func (e ValueNotFoundError) Error() string {
	return "Value not found in list: " + e.VarName
}
