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

type ForLoop struct {
	Init  int
	End   int
	Step  int
	State int
}

func (it *ForLoop) Boolean() Bool {
	return true
}

func (it *ForLoop) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoop) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoop) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ForLoop) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *ForLoop) Equals(Value) Bool {
	return false
}

func (it *ForLoop) IsIterable() Bool {
	return false
}

func (it *ForLoop) GetIterator() Value {
	return NilValue
}

func (it ForLoop) String() string {
	return fmt.Sprintf("ForLoop [i = %v, e = %v, d = %v, s = %v]", it.Init, it.End, it.Step, it.State)
}

func (it *ForLoop) Type() string {
	return "ForLoop"
}

type IForLoop struct {
	Iter  int
	State int
	Key   int
	Value int
}

func (it *IForLoop) Boolean() Bool {
	return true
}

func (it *IForLoop) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *IForLoop) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *IForLoop) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *IForLoop) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *IForLoop) Equals(Value) Bool {
	return false
}

func (it *IForLoop) IsIterable() Bool {
	return false
}

func (it *IForLoop) GetIterator() Value {
	return NilValue
}

func (it IForLoop) String() string {
	return fmt.Sprintf("IForLoop [i = %v, s = %v, k = %v, v = %v]", it.Iter, it.State, it.Key, it.Value)
}

func (it *IForLoop) Type() string {
	return "IForLoop"
}
