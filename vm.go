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
const stackSize = 256

type frame struct {
	code  []byte
	stack []Value
	ip    int
	ret   int
	op    byte
}

type VM struct {
	Frames       [callStackSize]frame
	Stack        [stackSize]Value
	Prelude      map[string]Value
	Module       *Module
	CurrentFrame *frame
	fp           int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m, Prelude: loadPrelude()}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	vm.CurrentFrame = &vm.Frames[vm.fp]
	vm.CurrentFrame.code = vm.Module.Code
	vm.CurrentFrame.stack = vm.Stack[:]
	ip := 8
	for {
		op := vm.CurrentFrame.code[ip]
		ip++
		switch op {
		case setG:
			scope := vm.CurrentFrame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			vm.Module.Store[vm.Module.Konstants[to].(String).Value] = vm.valueFrom(scope, from)
		case setL:
			scope := vm.CurrentFrame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := vm.CurrentFrame.code[ip]
			ip++
			vm.CurrentFrame.stack[to] = vm.valueFrom(scope, from)
		case move:
			from := vm.CurrentFrame.code[ip]
			ip++
			to := vm.CurrentFrame.code[ip]
			ip++
			vm.CurrentFrame.stack[to] = vm.CurrentFrame.stack[from]
		case prefix:
			op := vm.CurrentFrame.code[ip]
			ip++
			scope := vm.CurrentFrame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := vm.CurrentFrame.code[ip]
			ip++
			val, err := vm.valueFrom(scope, from).Prefix(op)
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.CurrentFrame.stack[to] = val
		case binop:
			op := vm.CurrentFrame.code[ip]
			ip++
			scopeLHS := vm.CurrentFrame.code[ip]
			ip++
			scopeRHS := vm.CurrentFrame.code[ip]
			ip++
			fromLHS := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			fromRHS := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := vm.CurrentFrame.code[ip]
			ip++
			val, err := vm.valueFrom(scopeLHS, fromLHS).Binop(op, vm.valueFrom(scopeRHS, fromRHS))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.CurrentFrame.stack[to] = val
		case iGet:
			scopeIndexable := vm.CurrentFrame.code[ip]
			ip++
			scopeIndex := vm.CurrentFrame.code[ip]
			ip++
			fromIndexable := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			fromIndex := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := vm.CurrentFrame.code[ip]
			ip++
			val, err := vm.valueFrom(scopeIndexable, fromIndexable).IGet(vm.valueFrom(scopeIndex, fromIndex))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.CurrentFrame.stack[to] = val
		case list:
			length := vm.CurrentFrame.code[ip]
			ip++
			from := vm.CurrentFrame.code[ip]
			ip++
			to := vm.CurrentFrame.code[ip]
			ip++
			xs := make([]Value, length)
			for i := 0; i < int(length); i++ {
				xs[i] = vm.CurrentFrame.stack[from]
				from++
			}
			vm.CurrentFrame.stack[to] = &List{Value: xs}
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxErrMsg, math.MaxUint16)
		}
	}
}

func (vm *VM) valueFrom(scope byte, from uint16) Value {
	switch scope {
	case rKonst:
		return vm.Module.Konstants[from]
	case rLocal:
		return vm.CurrentFrame.stack[from]
	case rGlobal:
		if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
	case rPrelude:
		if v, defined := vm.Prelude[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
	default:
		return NilValue
	}
}

func checkISACompatibility(m *Module) error {
	if m.Code[4] == major {
		return nil
	}
	return verror.New(m.Name, "The module was compilated with another ABI version", verror.FileErrMsg, 0)
}
