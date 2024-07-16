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

type ListIterator struct {
	List []Value
	Init int
	End  int
}

func (it *ListIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *ListIterator) Key() Value {
	return Integer(it.Init)
}

func (it *ListIterator) Value() Value {
	return it.List[it.Init]
}

func (it *ListIterator) Boolean() Bool {
	return true
}

func (it *ListIterator) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ListIterator) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ListIterator) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ListIterator) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *ListIterator) Equals(Value) Bool {
	return false
}

func (it *ListIterator) IsIterable() Bool {
	return false
}

func (it *ListIterator) Iterator() Value {
	return NilValue
}

func (it ListIterator) String() string {
	return fmt.Sprintf("ListIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *ListIterator) Type() string {
	return "ListIter"
}

type DocIterator struct {
	Keys []string
	Doc  map[string]Value
	Init int
	End  int
}

func (it *DocIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *DocIterator) Key() Value {
	return String{Value: it.Keys[it.Init]}
}

func (it *DocIterator) Value() Value {
	return it.Doc[it.Keys[it.Init]]
}

func (it *DocIterator) Boolean() Bool {
	return true
}

func (it *DocIterator) Prefix(byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *DocIterator) Binop(byte, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *DocIterator) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *DocIterator) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *DocIterator) Equals(Value) Bool {
	return false
}

func (it *DocIterator) IsIterable() Bool {
	return false
}

func (it *DocIterator) Iterator() Value {
	return NilValue
}

func (it DocIterator) String() string {
	return fmt.Sprintf("DocIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *DocIterator) Type() string {
	return "DocIter"
}
