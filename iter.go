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

func (it *ListIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ListIterator) Binop(uint64, Value) (Value, error) {
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

func (it *ListIterator) IsCallable() Bool {
	return false
}

func (it *ListIterator) Iterator() Value {
	return NilValue
}

func (it ListIterator) String() string {
	return fmt.Sprintf("ListIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *ListIterator) Clone() Value {
	return it
}

func (it *ListIterator) Type() string {
	return "ListIter"
}

type ObjectIterator struct {
	Keys []string
	Doc  map[string]Value
	Init int
	End  int
}

func (it *ObjectIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *ObjectIterator) Key() Value {
	return String{Value: it.Keys[it.Init]}
}

func (it *ObjectIterator) Value() Value {
	return it.Doc[it.Keys[it.Init]]
}

func (it *ObjectIterator) Boolean() Bool {
	return true
}

func (it *ObjectIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ObjectIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ObjectIterator) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *ObjectIterator) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *ObjectIterator) Equals(Value) Bool {
	return false
}

func (it *ObjectIterator) IsIterable() Bool {
	return false
}

func (it *ObjectIterator) IsCallable() Bool {
	return false
}

func (it *ObjectIterator) Iterator() Value {
	return NilValue
}

func (it ObjectIterator) String() string {
	return fmt.Sprintf("DocIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *ObjectIterator) Clone() Value {
	return it
}

func (it *ObjectIterator) Type() string {
	return "ObjIter"
}

type StringIterator struct {
	Runes []rune
	Init  int
	End   int
}

func (it *StringIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *StringIterator) Key() Value {
	return Integer(it.Init)
}

func (it *StringIterator) Value() Value {
	return String{Value: string(it.Runes[it.Init])}
}

func (it *StringIterator) Boolean() Bool {
	return true
}

func (it *StringIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *StringIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *StringIterator) IGet(Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (it *StringIterator) ISet(Value, Value) error {
	return verror.RuntimeError
}

func (it *StringIterator) Equals(Value) Bool {
	return false
}

func (it *StringIterator) IsIterable() Bool {
	return false
}

func (it *StringIterator) IsCallable() Bool {
	return false
}

func (it *StringIterator) Iterator() Value {
	return NilValue
}

func (it StringIterator) String() string {
	return fmt.Sprintf("StrIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *StringIterator) Clone() Value {
	return it
}

func (it *StringIterator) Type() string {
	return "StrIter"
}
