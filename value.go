package vida

import "fmt"

type Value any

type Nil struct{}

func NewNil() Nil {
	return Nil{}
}

func (n Nil) String() string {
	return "nil"
}

type Module struct {
	Konstants []Value
	Code      []byte
	Store     map[string]Value
}

func NewModule() *Module {
	m := Module{
		Store: make(map[string]Value),
		Code:  make([]byte, 0, 128),
	}
	return &m
}

func (m Module) String() string {
	return fmt.Sprintf("Module (%p)", &m)
}

type Function struct {
	FreeVars []Value
	Arity    int
	First    int
	Last     int
}

func (f Function) String() string {
	return fmt.Sprintf("Function (%p)", &f)
}
