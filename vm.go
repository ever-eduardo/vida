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
	Module   *Module
	Frames   [callStackSize]frame
	Stack    [stackSize]Value
	TopStore map[string]Value
	fp       int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	ip := 8
	for {
		op := frame.code[ip]
		ip++
		switch op {
		case setks:
			flag := frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			to := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			if flag == refKns {
				vm.Module.Store[vm.Module.Konstants[to].(string)] = vm.Module.Konstants[from]
			} else {
				vm.Module.Store[vm.Module.Konstants[to].(string)] = vm.Module.Store[vm.Module.Konstants[from].(string)]
			}
		case loadAtom:
			dest := frame.code[ip]
			ip++
			atom := frame.code[ip]
			ip++
			switch atom {
			case atomTrue:
				frame.stack[dest] = true
			case atomFalse:
				frame.stack[dest] = false
			case atomNil:
				frame.stack[dest] = globalNil
			}
		case loadGlobal:
			addr := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			dest := frame.code[ip]
			ip++
			frame.stack[dest] = vm.Module.Store[vm.Module.Konstants[addr].(string)]
			ip++
		case setGlobal:
			src := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			dest := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			if val, isPresent := vm.Module.Store[vm.Module.Konstants[src].(string)]; isPresent {
				vm.Module.Store[vm.Module.Konstants[dest].(string)] = val
			} else {
				vm.Module.Store[vm.Module.Konstants[dest].(string)] = globalNil
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
