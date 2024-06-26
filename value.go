package vida

import "fmt"

type Value any

type Nil struct{}

func (n Nil) String() string {
	return "nil"
}

type Module struct {
	Konstants []Value
	Code      []byte
	Name      string
	Store     map[string]Value
}

func NewModule(name string) Module {
	m := Module{
		Store:     make(map[string]Value),
		Code:      make([]byte, 0, 128),
		Konstants: make([]Value, 0, 16),
		Name:      name,
	}
	return m
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

type Closure struct {
	Function Function
	FreeVars []Value
}

func (c Closure) String() string {
	return "Function"
}
