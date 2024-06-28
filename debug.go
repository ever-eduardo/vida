package vida

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/verror"
)

func (vm *VM) Inspect(ip int) {
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

func (vm *VM) Debug() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	ip := 8
	for {
		vm.Inspect(ip)
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
		case locks:
			flag := frame.code[ip]
			ip++
			dest := frame.code[ip]
			ip++
			src := binary.NativeEndian.Uint16(frame.code[ip:])
			if flag == refKns {
				frame.stack[dest] = vm.Module.Konstants[src]
			} else {
				frame.stack[dest] = vm.Module.Store[vm.Module.Konstants[src].(string)]
			}
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxError, math.MaxUint16)
		}
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
