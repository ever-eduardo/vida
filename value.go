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
	Prefix(byte) (Value, error)
	Binop(byte, Value) (Value, error)
	IGet(Value) (Value, error)
	ISet(Value, Value) error
	Equals(Value) Bool
	IsIterable() Bool
	Iterator() Value
	String() string
	Type() string
}

type Nil struct{}

func (n Nil) Boolean() Bool {
	return Bool(false)
}

func (n Nil) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(true), nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (n Nil) Binop(op byte, rhs Value) (Value, error) {
	switch op {
	case byte(token.AND):
		return NilValue, nil
	case byte(token.OR):
		return rhs, nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (n Nil) IGet(index Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (n Nil) ISet(index, val Value) error {
	return verror.RuntimeError
}

func (n Nil) Equals(other Value) Bool {
	_, ok := other.(Nil)
	return Bool(ok)
}

func (n Nil) IsIterable() Bool {
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

type Bool bool

func (b Bool) Boolean() Bool {
	return b
}

func (b Bool) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return !b, nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (b Bool) Binop(op byte, rhs Value) (Value, error) {
	switch op {
	case byte(token.AND):
		if b {
			return rhs, nil
		}
		return b, nil
	case byte(token.OR):
		if b {
			return b, nil
		}
		return rhs, nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (b Bool) IGet(index Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (b Bool) ISet(index, val Value) error {
	return verror.RuntimeError
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

type String struct {
	Value string
}

func (s String) Boolean() Bool {
	return Bool(true)
}

func (s String) Binop(op byte, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case String:
		switch op {
		case byte(token.ADD):
			str := String{Value: s.Value + r.Value}
			return str, nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return s, nil
		case byte(token.LT):
			return Bool(s.Value < r.Value), nil
		case byte(token.LE):
			return Bool(s.Value <= r.Value), nil
		case byte(token.GT):
			return Bool(s.Value > r.Value), nil
		case byte(token.GE):
			return Bool(s.Value >= r.Value), nil
		}
	default:
		switch op {
		case byte(token.OR):
			return s, nil
		case byte(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (s String) IGet(index Value) (Value, error) {
	switch r := index.(type) {
	case Integer:
		l := Integer(len(s.Value))
		if r < 0 {
			r += l
		}
		if 0 <= r && r < l {
			return String{Value: s.Value[r : r+Integer(1)]}, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (s String) ISet(index, val Value) error {
	return verror.RuntimeError
}

func (s String) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.RuntimeError
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

func (s String) Iterator() Value {
	if r, ok := strToRunesMap[s.Value]; ok {
		return &StringIterator{Runes: r, Init: -1, End: len(r)}
	}
	r := []rune(s.Value)
	strToRunesMap[s.Value] = r
	return &StringIterator{Runes: r, Init: -1, End: len(r)}
}

func (s String) String() string {
	return s.Value
}

func (s String) Type() string {
	return "string"
}

type Integer int64

func (i Integer) Boolean() Bool {
	return Bool(true)
}

func (i Integer) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.SUB):
		return -i, nil
	case byte(token.NOT):
		return Bool(false), nil
	case byte(token.ADD):
		return i, nil
	}
	return NilValue, verror.RuntimeError
}

func (l Integer) Binop(op byte, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case Integer:
		switch op {
		case byte(token.ADD):
			return l + r, nil
		case byte(token.SUB):
			return l - r, nil
		case byte(token.MUL):
			return l * r, nil
		case byte(token.DIV):
			if r == 0 {
				return NilValue, verror.RuntimeError
			}
			return l / r, nil
		case byte(token.REM):
			if r == 0 {
				return NilValue, verror.RuntimeError
			}
			return l % r, nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return l, nil
		case byte(token.LT):
			return Bool(l < r), nil
		case byte(token.LE):
			return Bool(l <= r), nil
		case byte(token.GT):
			return Bool(l > r), nil
		case byte(token.GE):
			return Bool(l >= r), nil
		}
	case Float:
		switch op {
		case byte(token.ADD):
			return Float(Float(l) + r), nil
		case byte(token.SUB):
			return Float(Float(l) - r), nil
		case byte(token.MUL):
			return Float(Float(l) * r), nil
		case byte(token.DIV):
			return Float(Float(l) / r), nil
		case byte(token.REM):
			return Float(math.Remainder(float64(l), float64(r))), nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return l, nil
		case byte(token.LT):
			return Bool(Float(l) < r), nil
		case byte(token.LE):
			return Bool(Float(l) <= r), nil
		case byte(token.GT):
			return Bool(Float(l) > r), nil
		case byte(token.GE):
			return Bool(Float(l) >= r), nil
		}
	default:
		switch op {
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return l, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (i Integer) IGet(index Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (i Integer) ISet(index, val Value) error {
	return verror.RuntimeError
}

func (i Integer) Equals(other Value) Bool {
	if val, ok := other.(Integer); ok {
		return i == val
	}
	return false
}

func (i Integer) IsIterable() Bool {
	return false
}

func (i Integer) Iterator() Value {
	return NilValue
}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) Type() string {
	return "int"
}

type Float float64

func (f Float) Boolean() Bool {
	return Bool(true)
}

func (f Float) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.SUB):
		return -f, nil
	case byte(token.NOT):
		return Bool(false), nil
	case byte(token.ADD):
		return f, nil
	}
	return NilValue, verror.RuntimeError
}

func (f Float) Binop(op byte, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case Float:
		switch op {
		case byte(token.ADD):
			return f + r, nil
		case byte(token.SUB):
			return f - r, nil
		case byte(token.MUL):
			return f * r, nil
		case byte(token.DIV):
			return f / r, nil
		case byte(token.REM):
			return Float(math.Remainder(float64(f), float64(r))), nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return f, nil
		case byte(token.LT):
			return Bool(f < r), nil
		case byte(token.LE):
			return Bool(f <= r), nil
		case byte(token.GT):
			return Bool(f > r), nil
		case byte(token.GE):
			return Bool(f >= r), nil
		}
	case Integer:
		switch op {
		case byte(token.ADD):
			return f + Float(r), nil
		case byte(token.SUB):
			return f - Float(r), nil
		case byte(token.MUL):
			return f * Float(r), nil
		case byte(token.DIV):
			return f / Float(r), nil
		case byte(token.REM):
			return Float(math.Remainder(float64(f), float64(r))), nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return f, nil
		case byte(token.LT):
			return Bool(f < Float(r)), nil
		case byte(token.LE):
			return Bool(f <= Float(r)), nil
		case byte(token.GT):
			return Bool(f > Float(r)), nil
		case byte(token.GE):
			return Bool(f >= Float(r)), nil
		}
	default:
		switch op {
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return f, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (f Float) IGet(index Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (f Float) ISet(index, val Value) error {
	return verror.RuntimeError
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

func (f Float) Iterator() Value {
	return NilValue
}

func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f Float) Type() string {
	return "float"
}

type List struct {
	Value []Value
}

func (xs *List) Boolean() Bool {
	return Bool(true)
}

func (xs *List) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (xs *List) Binop(op byte, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case *List:
		switch op {
		case byte(token.ADD):
			rLen := len(r.Value)
			if rLen == 0 {
				return xs, nil
			}
			lLen := len(xs.Value)
			values := make([]Value, lLen+rLen)
			copy(values[:lLen], xs.Value)
			copy(values[lLen:], r.Value)
			return &List{Value: values}, nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return xs, nil
		}
	default:
		switch op {
		case byte(token.OR):
			return xs, nil
		case byte(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.RuntimeError
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
	return NilValue, verror.RuntimeError
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
	return verror.RuntimeError
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

type Document struct {
	Value map[string]Value
}

func (d *Document) Boolean() Bool {
	return Bool(true)
}

func (d *Document) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.RuntimeError
	}
}

func (d *Document) Binop(op byte, rhs Value) (Value, error) {
	switch r := rhs.(type) {
	case *Document:
		switch op {
		case byte(token.ADD):
			rLen := len(r.Value)
			if rLen == 0 {
				return d, nil
			}
			pairs := make(map[string]Value)
			for k, v := range d.Value {
				pairs[k] = v
			}
			for k, v := range r.Value {
				pairs[k] = v
			}
			return &Document{Value: pairs}, nil
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return d, nil
		}
	default:
		switch op {
		case byte(token.OR):
			return d, nil
		case byte(token.AND):
			return r, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (d *Document) IGet(index Value) (Value, error) {
	switch r := index.(type) {
	case String:
		if val, ok := d.Value[r.Value]; ok {
			return val, nil
		}
		return NilValue, nil
	}
	return NilValue, verror.RuntimeError
}

func (d *Document) ISet(index, val Value) error {
	switch r := index.(type) {
	case String:
		d.Value[r.Value] = val
		return nil
	}
	return verror.RuntimeError
}

func (d *Document) Equals(other Value) Bool {
	if val, ok := other.(*Document); ok {
		return d == val
	}
	return false
}

func (d *Document) IsIterable() Bool {
	return true
}

func (d *Document) Iterator() Value {
	size := len(d.Value)
	keys := make([]string, 0, size)
	for k := range d.Value {
		keys = append(keys, k)
	}
	return &DocIterator{Doc: d.Value, Init: -1, End: size, Keys: keys}
}

func (d *Document) String() string {
	if len(d.Value) == 0 {
		return "{}"
	}
	var r []string
	for k, v := range d.Value {
		r = append(r, fmt.Sprintf("%v: %v", k, v.String()))
	}
	return fmt.Sprintf("{%v}", strings.Join(r, ", "))
}

func (d *Document) Type() string {
	return "document"
}

type Module struct {
	Konstants []Value
	Code      []byte
	Name      string
	Store     map[string]Value
}

func newModule(name string) *Module {
	m := Module{
		Store:     make(map[string]Value),
		Code:      make([]byte, 0, 128),
		Konstants: make([]Value, 0, 32),
		Name:      name,
	}
	return &m
}

func (m Module) String() string {
	return fmt.Sprintf("Module <%v/>", m.Name)
}

type Function struct {
	FreeVarsCount int
	Arity         int
	First         int
	Last          int
}

func newFunction() *Function {
	return &Function{}
}

type Closure struct {
	Function *Function
	FreeVars []Value
}

func (c Closure) String() string {
	return "Function"
}

type GoFn func(args ...Value) (Value, error)

func (gfn GoFn) Boolean() Bool {
	return Bool(true)
}

func (gfn GoFn) Prefix(op byte) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (gfn GoFn) Binop(op byte, rhs Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (gfn GoFn) IGet(index Value) (Value, error) {
	return NilValue, verror.RuntimeError
}

func (gfn GoFn) ISet(index, val Value) error {
	return verror.RuntimeError
}

func (gfn GoFn) Equals(other Value) Bool {
	return false
}

func (gfn GoFn) IsIterable() Bool {
	return false
}

func (gfn GoFn) Iterator() Value {
	return NilValue
}

func (gfn GoFn) String() string {
	return "GFunction"
}

func (gfn GoFn) Type() string {
	return "function"
}
