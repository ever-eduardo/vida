package vida

import (
	"fmt"
)

type Script struct {
	Store        *[]Value
	Konstants    *[]Value
	MainFunction *Function
}

func newMainScript(name string) *Script {
	store := new([]Value)
	loadCoreLib(store)
	s := Script{
		Konstants:    nil,
		Store:        store,
		MainFunction: &Function{CoreFn: &CoreFunction{ScriptName: name}},
	}
	return &s
}

func newScript(name string, store *[]Value) *Script {
	loadCoreLib(store)
	s := Script{
		Konstants:    nil,
		Store:        store,
		MainFunction: &Function{CoreFn: &CoreFunction{ScriptName: name}},
	}
	return &s
}

func (s Script) String() string {
	return fmt.Sprintf("Script [%v]", s.MainFunction.CoreFn.ScriptName)
}
