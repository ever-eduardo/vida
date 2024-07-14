package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/verror"
)

type Iterator interface {
	Next() bool
	Key() Value
	Value() Value
}

type ForLoopRegisters struct {
	Init  int
	End   int
	Step  int
	State int
}

func (it *ForLoopRegisters) Boolean() Bool {
	return true
}

func (it *ForLoopRegisters) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopRegisters) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopRegisters) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopRegisters) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *ForLoopRegisters) Equals(Value) Bool {
	return false
}

func (it ForLoopRegisters) String() string {
	return fmt.Sprintf("ForLoopRegisters [i = %v, e = %v, d = %v, s = %v]", it.Init, it.End, it.Step, it.State)
}

func (it *ForLoopRegisters) Type() string {
	return "ForLoopRegisters"
}
