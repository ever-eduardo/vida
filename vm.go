package vida

import "encoding/binary"

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
	Module Module
	Frames [callStackSize]frame
	Stack  [stackSize]Value
	fp     int
}

func NewVM(m Module) VM {
	return VM{Module: m}
}

func (vm VM) Run() Result {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	ip := 8
	for {
		op := frame.code[ip]
		switch op {
		case setAtom:
			ip++
			addr := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			atom := frame.code[ip]
			ip++
			switch atom {
			case atomTrue:
				vm.Module.Store[vm.Module.Konstants[addr].(string)] = true
			case atomFalse:
				vm.Module.Store[vm.Module.Konstants[addr].(string)] = false
			case atomNil:
				vm.Module.Store[vm.Module.Konstants[addr].(string)] = globalNil
			}
		case loadAtom:
			ip++
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
			ip++
			addr := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			dest := frame.code[ip]
			ip++
			frame.stack[dest] = vm.Module.Store[vm.Module.Konstants[addr].(string)]
			ip++
		case setGlobal:
			ip++
			src := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			dest := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			if val, isPresent := vm.Module.Store[vm.Module.Konstants[src].(string)]; isPresent {
				vm.Module.Store[vm.Module.Konstants[dest].(string)] = val
			} else {
				vm.Module.Store[vm.Module.Konstants[dest].(string)] = globalNil
			}
		case stopRun:
			return Success
		default:
			return Failure
		}
	}
}
