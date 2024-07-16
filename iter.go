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
	List  *List
	Init  int
	End   int
	State int
}

func (it *ListIterator) Next() bool {
	ok := it.State < it.End
	if ok {
		it.Init = it.State
		it.State++
	}
	return ok
}

func (it *ListIterator) Key() Value {
	return Integer(it.Init)
}

func (it *ListIterator) Value() Value {
	return it.List.Value[it.Init]
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

func (it *ListIterator) GetIterator() Value {
	return NilValue
}

func (it ListIterator) String() string {
	return fmt.Sprintf("ListIter [i = %v, e = %v, s = %v]", it.Init, it.End, it.State)
}

func (it *ListIterator) Type() string {
	return "ListIter"
}

type DocIterator struct {
	Doc   *Document
	Keys  []string
	Init  int
	End   int
	State int
}

func (it *DocIterator) Next() bool {
	ok := it.State < it.End
	if ok {
		it.Init = it.State
		it.State++
	}
	return ok
}

func (it *DocIterator) Key() Value {
	return String{Value: it.Keys[it.Init]}
}

func (it *DocIterator) Value() Value {
	return it.Doc.Value[it.Keys[it.Init]]
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

func (it *DocIterator) GetIterator() Value {
	return NilValue
}

func (it DocIterator) String() string {
	return fmt.Sprintf("DocIter [i = %v, e = %v, s = %v]", it.Init, it.End, it.State)
}

func (it *DocIterator) Type() string {
	return "DocIter"
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
