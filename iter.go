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
	ReferenceSemanticsImpl
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
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ListIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ListIterator) IGet(Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ListIterator) ISet(Value, Value) error {
	return verror.ErrOpNotDefinedForIterators
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
	ReferenceSemanticsImpl
	Keys []string
	Obj  map[string]Value
	Init int
	End  int
}

func (it *ObjectIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *ObjectIterator) Key() Value {
	return &String{Value: it.Keys[it.Init]}
}

func (it *ObjectIterator) Value() Value {
	return it.Obj[it.Keys[it.Init]]
}

func (it *ObjectIterator) Boolean() Bool {
	return true
}

func (it *ObjectIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ObjectIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ObjectIterator) IGet(Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *ObjectIterator) ISet(Value, Value) error {
	return verror.ErrOpNotDefinedForIterators
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
	ReferenceSemanticsImpl
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
	return &String{Value: string(it.Runes[it.Init]), Runes: it.Runes[it.Init : it.Init+1]}
}

func (it *StringIterator) Boolean() Bool {
	return true
}

func (it *StringIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *StringIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *StringIterator) IGet(Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *StringIterator) ISet(Value, Value) error {
	return verror.ErrOpNotDefinedForIterators
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

type IntegerIterator struct {
	ReferenceSemanticsImpl
	Init Integer
	End  Integer
}

func (it *IntegerIterator) Next() bool {
	it.Init++
	return it.Init < it.End
}

func (it *IntegerIterator) Key() Value {
	return it.Init
}

func (it *IntegerIterator) Value() Value {
	return it.Init
}

func (it *IntegerIterator) Boolean() Bool {
	return true
}

func (it *IntegerIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *IntegerIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *IntegerIterator) IGet(Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (it *IntegerIterator) ISet(Value, Value) error {
	return verror.ErrOpNotDefinedForIterators
}

func (it *IntegerIterator) Equals(Value) Bool {
	return false
}

func (it *IntegerIterator) IsIterable() Bool {
	return false
}

func (it *IntegerIterator) IsCallable() Bool {
	return false
}

func (it *IntegerIterator) Iterator() Value {
	return NilValue
}

func (it IntegerIterator) String() string {
	return fmt.Sprintf("IntIter [i = %v, e = %v]", it.Init, it.End)
}

func (it *IntegerIterator) Clone() Value {
	return it
}

func (it *IntegerIterator) Type() string {
	return "IntIter"
}

type BytesIterator struct {
	ReferenceSemanticsImpl
	Bytes []byte
	Init  int
	End   int
}

func (bi *BytesIterator) Next() bool {
	bi.Init++
	return bi.Init < bi.End
}

func (bi *BytesIterator) Key() Value {
	return Integer(bi.Init)
}

func (bi *BytesIterator) Value() Value {
	return Integer(bi.Bytes[bi.Init])
}

func (bi *BytesIterator) Boolean() Bool {
	return true
}

func (bi *BytesIterator) Prefix(uint64) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (bi *BytesIterator) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (bi *BytesIterator) IGet(Value) (Value, error) {
	return NilValue, verror.ErrOpNotDefinedForIterators
}

func (bi *BytesIterator) ISet(Value, Value) error {
	return verror.ErrOpNotDefinedForIterators
}

func (bi *BytesIterator) Equals(Value) Bool {
	return false
}

func (bi *BytesIterator) IsIterable() Bool {
	return false
}

func (bi *BytesIterator) IsCallable() Bool {
	return false
}

func (bi *BytesIterator) Iterator() Value {
	return NilValue
}

func (bi BytesIterator) String() string {
	return fmt.Sprintf("BytesIter [i = %v, e = %v]", bi.Init, bi.End)
}

func (bi *BytesIterator) Clone() Value {
	return bi
}

func (bi *BytesIterator) Type() string {
	return "BytesIter"
}
