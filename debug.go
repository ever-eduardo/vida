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
			message := fmt.Sprintf("Unknown instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxError, math.MaxUint16)
		}
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
