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
	fmt.Print("[")
	for _, v := range vm.Stack {
		if v == nil {
			fmt.Printf(" %5v ", "_")
		} else {
			fmt.Printf(" %5v ", v)
		}
	}
	fmt.Print("]")
	fmt.Println()
	fmt.Printf("Press 'Enter' to continue => ")
	fmt.Scanf(" ")
}

func (vm *VM) Debug() (Result, error) {
	frame := &vm.Frames[vm.fp]
	frame.code = vm.Module.Code
	frame.stack = vm.Stack[:]
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
			switch flag {
			case refKns:
				vm.Module.Store[vm.Module.Konstants[to].(string)] = vm.Module.Konstants[from]
			case refStr:
				if v, defined := vm.Module.Store[vm.Module.Konstants[from].(string)]; defined {
					vm.Module.Store[vm.Module.Konstants[to].(string)] = v
				} else {
					vm.Module.Store[vm.Module.Konstants[to].(string)] = globalNil
				}
			case refReg:
				vm.Module.Store[vm.Module.Konstants[to].(string)] = frame.stack[from]
			}
		case locks:
			flag := frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(frame.code[ip:])
			ip += 2
			to := frame.code[ip]
			ip++
			if flag == refKns {
				frame.stack[to] = vm.Module.Konstants[from]
			} else {
				if v, defined := vm.Module.Store[vm.Module.Konstants[from].(string)]; defined {
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
