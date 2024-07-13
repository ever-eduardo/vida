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

type ForLoopState struct {
	Init  int
	End   int
	Step  int
	State int
}

func (it *ForLoopState) Boolean() Bool {
	return true
}

func (it *ForLoopState) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopState) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopState) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoopState) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *ForLoopState) Equals(Value) Bool {
	return false
}

func (it ForLoopState) String() string {
	return fmt.Sprintf("ForLoopState [i = %v, e = %v, d = %v, s = %v]", it.Init, it.End, it.Step, it.State)
}

func (it *ForLoopState) Type() string {
	return "ForLoopState"
}
