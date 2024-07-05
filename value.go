package vida

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Value interface {
	Boolean() Bool
	Prefix(byte) (Value, error)
	Binop(byte, Value) (Value, error)
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
		return NilValue, errors.New(verror.RunTimeError)
	}
}

func (n Nil) Binop(op byte, rhs Value) (Value, error) {
	switch op {
	case byte(token.AND):
		return NilValue, nil
	case byte(token.OR):
		return rhs, nil
	default:
		return NilValue, errors.New(verror.RunTimeError)
	}
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
		return NilValue, errors.New(verror.RunTimeError)
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
		return NilValue, errors.New(verror.RunTimeError)
	}
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
	switch op {
	case byte(token.OR):
		return s.Boolean() || rhs.Boolean(), nil
	case byte(token.AND):
		return s.Boolean() && rhs.Boolean(), nil
	default:
		return NilValue, errors.New(verror.RunTimeError)
	}
}

func (s String) Prefix(op byte) (Value, error) {
	switch op {
	case byte(token.NOT):
		return Bool(len(s.Value) != 0), nil
	default:
		return NilValue, errors.New(verror.RunTimeError)
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
	case byte(token.ADD):
		return i, nil
	case byte(token.NOT):
		return Bool(false), nil
	}
	return NilValue, errors.New(verror.RunTimeError)
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
				return NilValue, errors.New(verror.RunTimeError)
			}
			return l / r, nil
		case byte(token.REM):
			if r == 0 {
				return NilValue, errors.New(verror.RunTimeError)
			}
			return l % r, nil
		default:
			switch op {
			case byte(token.AND):
				return r, nil
			case byte(token.OR):
				return l, nil
			}
		}
	default:
		switch op {
		case byte(token.AND):
			return r, nil
		case byte(token.OR):
			return l, nil
		}
	}
	return NilValue, errors.New(verror.RunTimeError)
}

func (i Integer) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Integer) Type() string {
	return "int"
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

func (gfn GoFn) String() string {
	return "Function"
}
