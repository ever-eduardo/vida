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
		l := len(s.Value)
		if -l <= int(r) && int(r) <= l-1 {
			if r < 0 {
				return String{Value: s.Value[l+int(r) : l+int(r)+1]}, nil
			}
			return String{Value: s.Value[int(r) : int(r)+1]}, nil
		}
	}
	return NilValue, verror.RuntimeError
}

func (s String) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.RuntimeError
	}
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
		l := len(xs.Value)
		if -l <= int(r) && int(r) <= l-1 {
			if r < 0 {
				return xs.Value[l+int(r)], nil
			}
			return xs.Value[r], nil
		}
	}
	return NilValue, verror.RuntimeError
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

func (gfn GoFn) String() string {
	return "GFunction"
}

func (gfn GoFn) Type() string {
	return "function"
}
