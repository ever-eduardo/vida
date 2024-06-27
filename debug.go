package vida

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/verror"
)

func (vm VM) Inspect(ip int) {
	clear()
	fmt.Println("Running", vm.Module.Name)
	fmt.Println("Store")
	for k, v := range vm.Module.Store {
		fmt.Println(k, " : ", v)
	}
	fmt.Println("Konst")
	for i, v := range vm.Module.Konstants {
		fmt.Println(i, " : ", v)
	}
	fmt.Println("Code")
	for i, v := range vm.Module.Code {
		if i == ip {
			fmt.Printf("[%v : %v], ", i, v)
		} else {
			fmt.Printf("%v, ", v)
		}
	}
	fmt.Println()
	fmt.Println("Stack")
	for _, v := range vm.Stack {
		fmt.Printf("%v, ", v)
	}
	fmt.Println()
	fmt.Printf("Press 'Enter' to continue => ")
	fmt.Scanf(" ")
}

func (vm VM) Debug() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	ip := 8
	for {
		vm.Inspect(ip)
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
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxError, math.MaxUint16)
		}
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
