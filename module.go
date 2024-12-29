package vida

import (
	"fmt"
)

type Module struct {
	Store        *[]Value
	Konstants    *[]Value
	MainFunction *Function
}

func newMainModule(name string) *Module {
	store := new([]Value)
	loadCoreLib(store)
	m := Module{
		Konstants:    nil,
		Store:        store,
		MainFunction: &Function{CoreFn: &CoreFunction{ModuleName: name}},
	}
	return &m
}

func newSubModule(name string, store *[]Value) *Module {
	loadCoreLib(store)
	m := Module{
		Konstants:    nil,
		Store:        store,
		MainFunction: &Function{CoreFn: &CoreFunction{ModuleName: name}},
	}
	return &m
}

func (m Module) String() string {
	return fmt.Sprintf("Module [%v]", m.MainFunction.CoreFn.ModuleName)
}
