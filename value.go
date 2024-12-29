package vida

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Value interface {
	Boolean() Bool
	Prefix(uint64) (Value, error)
	Binop(uint64, Value) (Value, error)
	IGet(Value) (Value, error)
	ISet(Value, Value) error
	Equals(Value) Bool
	IsIterable() Bool
	Iterator() Value
	IsCallable() Bool
	String() string
	Type() string
	Clone() Value
}

type Nil struct{}

func (n Nil) Boolean() Bool {
	return Bool(false)
}

func (n Nil) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(true), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (n Nil) Binop(op uint64, rhs Value) (Value, error) {
	switch op {
	case uint64(token.AND):
		return NilValue, nil
	case uint64(token.OR):
		return rhs, nil
	default:
		return NilValue, verror.ErrBinaryOpNotDefined
	}
}

func (n Nil) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (n Nil) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (n Nil) Equals(other Value) Bool {
	_, ok := other.(Nil)
	return Bool(ok)
}

func (n Nil) IsIterable() Bool {
	return false
}

func (n Nil) IsCallable() Bool {
	return false
}

func (n Nil) Iterator() Value {
	return NilValue
}

func (n Nil) String() string {
	return "nil"
}

func (n Nil) Type() string {
	return "nil"
}

func (n Nil) Clone() Value {
	return n
}

type Bool bool

func (b Bool) Boolean() Bool {
	return b
}

func (b Bool) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return !b, nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (b Bool) Binop(op uint64, rhs Value) (Value, error) {
	switch op {
	case uint64(token.AND):
		if b {
			return rhs, nil
		}
		return b, nil
	case uint64(token.OR):
		if b {
			return b, nil
		}
		return rhs, nil
	default:
		return NilValue, verror.ErrBinaryOpNotDefined
	}
}

func (b Bool) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (b Bool) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (b Bool) Equals(other Value) Bool {
	if val, ok := other.(Bool); ok {
		return b == val
	}
	return false
}

func (b Bool) IsIterable() Bool {
	return false
}

func (b Bool) IsCallable() Bool {
	return false
}

func (b Bool) Iterator() Value {
	return NilValue
}

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Bool) Type() string {
	return "bool"
}

func (b Bool) Clone() Value {
	return b
}

type String struct {
	Runes []rune
	Value string
}

func (s String) Boolean() Bool {
	return Bool(true)
}

func (s String) Binop(op uint64, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case String:
		switch op {
		case uint64(token.ADD):
			str := String{Value: s.Value + r.Value}
			return str, nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return s, nil
		case uint64(token.LT):
			return Bool(s.Value < r.Value), nil
		case uint64(token.LE):
			return Bool(s.Value <= r.Value), nil
		case uint64(token.GT):
			return Bool(s.Value > r.Value), nil
		case uint64(token.GE):
			return Bool(s.Value >= r.Value), nil
		}
	default:
		switch op {
		case uint64(token.OR):
			return s, nil
		case uint64(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (s String) IGet(index Value) (Value, error) {
	switch r := index.(type) {
	case Integer:
		if s.Runes == nil {
			s.Runes = []rune(s.Value)
		}
		l := Integer(len(s.Runes))
		if r < 0 {
			r += l
		}
		if 0 <= r && r < l {
			sr := s.Runes[r : r+Integer(1)]
			return String{Value: string(sr), Runes: sr}, nil
		}
	}
	return NilValue, verror.ErrValueNotIndexable
}

func (s String) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (s String) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (s String) Equals(other Value) Bool {
	if val, ok := other.(String); ok {
		return s.Value == val.Value
	}
	return false
}

func (s String) IsIterable() Bool {
	return true
}

func (s String) IsCallable() Bool {
	return false
}

func (s String) Iterator() Value {
	if s.Runes == nil {
		s.Runes = []rune(s.Value)
	}
	return &StringIterator{Runes: s.Runes, Init: -1, End: len(s.Runes)}
}

func (s String) String() string {
	return s.Value
}

func (s String) Type() string {
	return "string"
}

func (s String) Clone() Value {
	return &String{Runes: s.Runes, Value: s.Value}
}

type Integer int64

func (i Integer) Boolean() Bool {
	return Bool(true)
}

func (i Integer) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.SUB):
		return -i, nil
	case uint64(token.NOT):
		return Bool(false), nil
	case uint64(token.ADD):
		return i, nil
	case uint64(token.TILDE):
		return Integer(^uint32(i)), nil
	}
	return NilValue, verror.ErrPrefixOpNotDefined
}

func (l Integer) Binop(op uint64, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case Integer:
		switch op {
		case uint64(token.ADD):
			return l + r, nil
		case uint64(token.SUB):
			return l - r, nil
		case uint64(token.MUL):
			return l * r, nil
		case uint64(token.DIV):
			if r == 0 {
				return NilValue, verror.ErrDivisionByZero
			}
			return l / r, nil
		case uint64(token.REM):
			if r == 0 {
				return NilValue, verror.ErrDivisionByZero
			}
			return l % r, nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return l, nil
		case uint64(token.LT):
			return Bool(l < r), nil
		case uint64(token.LE):
			return Bool(l <= r), nil
		case uint64(token.GT):
			return Bool(l > r), nil
		case uint64(token.GE):
			return Bool(l >= r), nil
		case uint64(token.BXOR):
			return Integer(uint32(l) ^ uint32(r)), nil
		case uint64(token.BOR):
			return Integer(uint32(l) | uint32(r)), nil
		case uint64(token.BAND):
			return Integer(uint32(l) & uint32(r)), nil
		case uint64(token.BSHL):
			return Integer(uint32(l) << uint32(r)), nil
		case uint64(token.BSHR):
			return Integer(uint32(l) >> uint32(r)), nil
		}
	case Float:
		switch op {
		case uint64(token.ADD):
			return Float(Float(l) + r), nil
		case uint64(token.SUB):
			return Float(Float(l) - r), nil
		case uint64(token.MUL):
			return Float(Float(l) * r), nil
		case uint64(token.DIV):
			return Float(Float(l) / r), nil
		case uint64(token.REM):
			return Float(math.Remainder(float64(l), float64(r))), nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return l, nil
		case uint64(token.LT):
			return Bool(Float(l) < r), nil
		case uint64(token.LE):
			return Bool(Float(l) <= r), nil
		case uint64(token.GT):
			return Bool(Float(l) > r), nil
		case uint64(token.GE):
			return Bool(Float(l) >= r), nil
		}
	default:
		switch op {
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return l, nil
		}
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (i Integer) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (i Integer) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (i Integer) Equals(other Value) Bool {
	if val, ok := other.(Integer); ok {
		return i == val
	}
	return false
}

func (i Integer) IsIterable() Bool {
	return true
}

func (i Integer) IsCallable() Bool {
	return false
}

func (i Integer) Iterator() Value {
	if i < 0 {
		i = -i
	}
	return &IntegerIterator{Init: -1, End: i}
}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) Type() string {
	return "int"
}

func (i Integer) Clone() Value {
	return i
}

type Float float64

func (f Float) Boolean() Bool {
	return Bool(true)
}

func (f Float) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.SUB):
		return -f, nil
	case uint64(token.NOT):
		return Bool(false), nil
	case uint64(token.ADD):
		return f, nil
	}
	return NilValue, verror.ErrPrefixOpNotDefined
}

func (f Float) Binop(op uint64, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case Float:
		switch op {
		case uint64(token.ADD):
			return f + r, nil
		case uint64(token.SUB):
			return f - r, nil
		case uint64(token.MUL):
			return f * r, nil
		case uint64(token.DIV):
			return f / r, nil
		case uint64(token.REM):
			return Float(math.Remainder(float64(f), float64(r))), nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return f, nil
		case uint64(token.LT):
			return Bool(f < r), nil
		case uint64(token.LE):
			return Bool(f <= r), nil
		case uint64(token.GT):
			return Bool(f > r), nil
		case uint64(token.GE):
			return Bool(f >= r), nil
		}
	case Integer:
		switch op {
		case uint64(token.ADD):
			return f + Float(r), nil
		case uint64(token.SUB):
			return f - Float(r), nil
		case uint64(token.MUL):
			return f * Float(r), nil
		case uint64(token.DIV):
			return f / Float(r), nil
		case uint64(token.REM):
			return Float(math.Remainder(float64(f), float64(r))), nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return f, nil
		case uint64(token.LT):
			return Bool(f < Float(r)), nil
		case uint64(token.LE):
			return Bool(f <= Float(r)), nil
		case uint64(token.GT):
			return Bool(f > Float(r)), nil
		case uint64(token.GE):
			return Bool(f >= Float(r)), nil
		}
	default:
		switch op {
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return f, nil
		}
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (f Float) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (f Float) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (f Float) Equals(other Value) Bool {
	if val, ok := other.(Float); ok {
		return f == val
	}
	return false
}

func (f Float) IsIterable() Bool {
	return false
}

func (f Float) IsCallable() Bool {
	return false
}

func (f Float) Iterator() Value {
	return NilValue
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f Float) Type() string {
	return "float"
}

func (f Float) Clone() Value {
	return f
}

type List struct {
	Value []Value
}

func (xs *List) Boolean() Bool {
	return Bool(true)
}

func (xs *List) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (xs *List) Binop(op uint64, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case *List:
		switch op {
		case uint64(token.ADD):
			rLen := len(r.Value)
			if rLen == 0 {
				return xs, nil
			}
			lLen := len(xs.Value)
			values := make([]Value, lLen+rLen)
			copy(values[:lLen], xs.Value)
			copy(values[lLen:], r.Value)
			return &List{Value: values}, nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return xs, nil
		}
	default:
		switch op {
		case uint64(token.OR):
			return xs, nil
		case uint64(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (xs *List) IGet(index Value) (Value, error) {
	switch r := index.(type) {
	case Integer:
		l := Integer(len(xs.Value))
		if r < 0 {
			r += l
		}
		if 0 <= r && r < l {
			return xs.Value[r], nil
		}
	}
	return NilValue, verror.ErrValueNotIndexable
}

func (xs *List) ISet(index, val Value) error {
	switch r := index.(type) {
	case Integer:
		l := Integer(len(xs.Value))
		if r < 0 {
			r += l
		}
		if 0 <= r && r < l {
			xs.Value[r] = val
			return nil
		}
	}
	return verror.ErrValueNotIndexable
}

func (xs *List) Equals(other Value) Bool {
	if val, ok := other.(*List); ok {
		return xs == val
	}
	return false
}

func (xs *List) IsIterable() Bool {
	return true
}

func (xs *List) IsCallable() Bool {
	return false
}

func (xs *List) Iterator() Value {
	return &ListIterator{List: xs.Value, Init: -1, End: len(xs.Value)}
}

func (xs List) String() string {
	if len(xs.Value) == 0 {
		return "[]"
	}
	var r []string
	for _, v := range xs.Value {
		r = append(r, v.String())
	}
	return fmt.Sprintf("[%v]", strings.Join(r, ", "))
}

func (xs *List) Type() string {
	return "list"
}

func (xs *List) Clone() Value {
	c := make([]Value, len(xs.Value))
	for i, v := range xs.Value {
		c[i] = v.Clone()
	}
	return &List{Value: c}
}

type Object struct {
	Keys  []string
	Value map[string]Value
}

func (o *Object) Boolean() Bool {
	return true
}

func (o *Object) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (o *Object) Binop(op uint64, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case *Object:
		switch op {
		case uint64(token.ADD):
			pairs := make(map[string]Value)
			keys := make([]string, 0)
			for k, v := range o.Value {
				pairs[k] = v
				keys = append(keys, k)
			}
			for k, v := range r.Value {
				if _, isPresent := pairs[k]; !isPresent {
					keys = append(keys, k)
				}
				pairs[k] = v
			}
			return &Object{Value: pairs, Keys: keys}, nil
		case uint64(token.AND):
			return r, nil
		case uint64(token.OR):
			return o, nil
		}
	default:
		switch op {
		case uint64(token.OR):
			return o, nil
		case uint64(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (o *Object) IGet(index Value) (Value, error) {
	if val, ok := o.Value[index.String()]; ok {
		return val, nil
	}
	return NilValue, nil
}

func (o *Object) ISet(index, val Value) error {
	k := index.String()
	if _, isPresent := o.Value[k]; !isPresent {
		o.Keys = append(o.Keys, k)
	}
	o.Value[k] = val
	return nil
}

func (o *Object) Equals(other Value) Bool {
	if val, ok := other.(*Object); ok {
		return o == val
	}
	return false
}

func (o *Object) IsIterable() Bool {
	return true
}

func (o *Object) IsCallable() Bool {
	return false
}

func (o *Object) Iterator() Value {
	return &ObjectIterator{Obj: o.Value, Init: -1, End: len(o.Value), Keys: o.Keys}
}

func (o *Object) String() string {
	if len(o.Value) == 0 {
		return "{}"
	}
	var r []string
	for _, v := range o.Keys {
		r = append(r, fmt.Sprintf("%v: %v", v, o.Value[v].String()))
	}
	return fmt.Sprintf("{%v}", strings.Join(r, ", "))
}

func (o *Object) Type() string {
	return "object"
}

func (o *Object) Clone() Value {
	m := make(map[string]Value)
	k := make([]string, len(o.Keys))
	copy(k, o.Keys)
	for k, v := range o.Value {
		m[k] = v.Clone()
	}
	return &Object{Value: m, Keys: k}
}

func (o *Object) UpdateKeys() {
	keys := make([]string, 0, len(o.Value))
	for k := range o.Value {
		keys = append(keys, k)
	}
	o.Keys = keys
}

type freeInfo struct {
	Index   int
	IsLocal Bool
	Id      string
}

type CoreFunction struct {
	Code       []uint64
	Info       []freeInfo
	Free       int
	Arity      int
	IsVar      bool
	ModuleName string
}

func (c *CoreFunction) Boolean() Bool {
	return true
}

func (c *CoreFunction) Prefix(uint64) (Value, error) {
	return NilValue, verror.ErrPrefixOpNotDefined
}

func (c *CoreFunction) Binop(uint64, Value) (Value, error) {
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (c *CoreFunction) IGet(Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (c *CoreFunction) ISet(Value, Value) error {
	return verror.ErrValueNotIndexable
}

func (c *CoreFunction) Equals(other Value) Bool {
	if f, ok := other.(*CoreFunction); ok {
		return c == f
	}
	return false
}

func (c *CoreFunction) IsIterable() Bool {
	return false
}

func (c *CoreFunction) IsCallable() Bool {
	return false
}

func (c *CoreFunction) Iterator() Value {
	return NilValue
}

func (c *CoreFunction) Type() string {
	return "corefunction"
}

func (f CoreFunction) String() string {
	return fmt.Sprintf("[a = %v, v = %v, f = %v, i = %v, mod = %v]", f.Arity, f.IsVar, f.Free, f.Info, f.ModuleName)
}

func (f *CoreFunction) Clone() Value {
	return f
}

type Function struct {
	Free   []Value
	CoreFn *CoreFunction
}

func (f *Function) Boolean() Bool {
	return true
}

func (f *Function) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (f *Function) Binop(op uint64, r Value) (Value, error) {
	switch op {
	case uint64(token.OR):
		return f, nil
	case uint64(token.AND):
		return r, nil
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (f *Function) IGet(Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (f *Function) ISet(Value, Value) error {
	return verror.ErrValueNotIndexable
}

func (f *Function) Equals(other Value) Bool {
	if o, ok := other.(*Function); ok {
		return f == o
	}
	return false
}

func (f *Function) IsIterable() Bool {
	return false
}

func (f *Function) IsCallable() Bool {
	return true
}

func (f *Function) Iterator() Value {
	return NilValue
}

func (f *Function) Type() string {
	return "function"
}

func (f *Function) Clone() Value {
	return f
}

func (f Function) String() string {
	return fmt.Sprintf("Function(%v, Free = %v)", f.CoreFn, f.Free)
}

type GFn func(args ...Value) (Value, error)

func (gfn GFn) Boolean() Bool {
	return Bool(true)
}

func (gfn GFn) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (gfn GFn) Binop(op uint64, r Value) (Value, error) {
	switch op {
	case uint64(token.OR):
		return gfn, nil
	case uint64(token.AND):
		return r, nil
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (gfn GFn) IGet(index Value) (Value, error) {
	return NilValue, verror.ErrValueNotIndexable
}

func (gfn GFn) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (gfn GFn) Equals(other Value) Bool {
	return false
}

func (gfn GFn) IsIterable() Bool {
	return false
}

func (gfn GFn) IsCallable() Bool {
	return true
}

func (gfn GFn) Iterator() Value {
	return NilValue
}

func (gfn GFn) String() string {
	return "GFn"
}

func (gFn GFn) Clone() Value {
	return gFn
}

func (gfn GFn) Type() string {
	return "function"
}

type Error struct {
	Message Value
}

func (e Error) Boolean() Bool {
	return false
}

func (e Error) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(true), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (e Error) Binop(op uint64, rhs Value) (Value, error) {
	switch op {
	case uint64(token.AND):
		return e, nil
	case uint64(token.OR):
		return rhs, nil
	default:
		return NilValue, verror.ErrBinaryOpNotDefined
	}
}

func (e Error) IGet(index Value) (Value, error) {
	if val, ok := index.(String); ok && val.Value == "message" {
		return e.Message, nil
	}
	return NilValue, nil
}

func (e Error) ISet(index, val Value) error {
	return verror.ErrValueNotIndexable
}

func (e Error) Equals(other Value) Bool {
	v, ok := other.(Error)
	return Bool(ok) && e.Message.Equals(v.Message)
}

func (e Error) IsIterable() Bool {
	return false
}

func (e Error) IsCallable() Bool {
	return false
}

func (e Error) Iterator() Value {
	return NilValue
}

func (e Error) String() string {
	return fmt.Sprintf("Error: %v", e.Message.String())
}

func (e Error) Type() string {
	return "error"
}

func (e Error) Clone() Value {
	return e
}
