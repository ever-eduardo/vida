package vida

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/verror"
)

type Result string

const Success Result = "Success"
const Failure Result = "Failure"

const callStackSize = 1024
const stackSize = 16

type frame struct {
	code  []byte
	stack []Value
	ip    int
	ret   int
	op    byte
}

type VM struct {
	Module  *Module
	Frames  [callStackSize]frame
	Stack   [stackSize]Value
	Prelude map[string]Value
	fp      int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m, Prelude: loadPrelude()}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	frame.stack = vm.Stack[:]
	ip := 8
	for {
		op := frame.code[ip]
		ip++
		switch op {
		case setG:
			scope := frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			to := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			switch scope {
			case rKonst:
				vm.Module.Store[vm.Module.Konstants[to].(String).Value] = vm.Module.Konstants[from]
			case rGlobal:
				if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
					vm.Module.Store[vm.Module.Konstants[to].(String).Value] = v
				} else {
					vm.Module.Store[vm.Module.Konstants[to].(String).Value] = globalNil
				}
			case rLocal:
				vm.Module.Store[vm.Module.Konstants[to].(String).Value] = frame.stack[from]
			case rPrelude:
				if v, defined := vm.Prelude[vm.Module.Konstants[from].(String).Value]; defined {
					vm.Module.Store[vm.Module.Konstants[to].(String).Value] = v
				} else {
					vm.Module.Store[vm.Module.Konstants[to].(String).Value] = globalNil
				}
			}
		case setL:
			scope := frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			to := frame.code[ip]
			ip++
			switch scope {
			case rKonst:
				frame.stack[to] = vm.Module.Konstants[from]
			case rLocal:
				frame.stack[to] = frame.stack[from]
			case rGlobal:
				if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
					frame.stack[to] = v
				} else {
					frame.stack[to] = globalNil
				}
			case rPrelude:
				if v, defined := vm.Prelude[vm.Module.Konstants[from].(String).Value]; defined {
					frame.stack[to] = v
				} else {
					frame.stack[to] = globalNil
				}
			}
		case move:
			from := frame.code[ip]
			ip++
			to := frame.code[ip]
			ip++
			frame.stack[to] = frame.stack[from]
		case not:
			scope := frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			to := frame.code[ip]
			ip++
			switch scope {
			case rKonst:
				frame.stack[to] = Bool(!vm.Module.Konstants[from].Boolean())
			case rLocal:
				frame.stack[to] = Bool(!frame.stack[from].Boolean())
			case rGlobal:
				if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
					frame.stack[to] = Bool(!v.Boolean())
				} else {
					frame.stack[to] = Bool(true)
				}
			case rPrelude:
				if v, defined := vm.Prelude[vm.Module.Konstants[from].(String).Value]; defined {
					frame.stack[to] = Bool(!v.Boolean())
				} else {
					frame.stack[to] = Bool(false)
				}
			}
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxError, math.MaxUint16)
		}
	}
}

func checkISACompatibility(m *Module) error {
	if m.Code[4] == major {
		return nil
	}
	return verror.New(m.Name, "The module was compilated with another ABI version", verror.FileError, 0)
}
