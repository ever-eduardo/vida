package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/token"
)

type Value interface {
	Prefix(byte) Value
	String() string
	Type() string
}

type Nil struct{}

func (n Nil) Prefix(op byte) Value {
	switch op {
	case byte(token.NOT):
		return Bool(true)
	default:
		return globalNil
	}
}

func (n Nil) String() string {
	return "nil"
}

func (n Nil) Type() string {
	return "nil"
}

type Bool bool

func (b Bool) Prefix(op byte) Value {
	switch op {
	case byte(token.NOT):
		return !b
	default:
		return globalNil
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

func (s String) Prefix(op byte) Value {
	switch op {
	case byte(token.NOT):
		return Bool(len(s.Value) != 0)
	default:
		return globalNil
	}
}

func (s String) String() string {
	return s.Value
}

func (s String) Type() string {
	return "string"
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
