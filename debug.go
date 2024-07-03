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
	vm.CurrentFrame = &vm.Frames[vm.fp]
	vm.CurrentFrame.code = vm.Module.Code
	vm.CurrentFrame.stack = vm.Stack[:]
	ip := 8
	for {
		vm.Inspect(ip)
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
			vm.CurrentFrame.stack[to] = vm.valueFrom(scope, from).Prefix(op)
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
			l := vm.valueFrom(scopeLHS, fromLHS)
			r := vm.valueFrom(scopeRHS, fromRHS)
			val, _ := l.Binary(op, r)
			vm.CurrentFrame.stack[to] = val
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
