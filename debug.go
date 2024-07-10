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
		case slice:
			mode := vm.CurrentFrame.code[ip]
			ip++
			scopeV := vm.CurrentFrame.code[ip]
			ip++
			scopeL := vm.CurrentFrame.code[ip]
			ip++
			scopeR := vm.CurrentFrame.code[ip]
			ip++
			fromV := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			fromL := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			fromR := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			to := vm.CurrentFrame.code[ip]
			ip++
			val, err := vm.processSlice(mode, fromV, fromL, fromR, scopeV, scopeL, scopeR)
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
		case record:
			length := vm.CurrentFrame.code[ip]
			ip++
			from := vm.CurrentFrame.code[ip]
			ip++
			to := vm.CurrentFrame.code[ip]
			ip++
			rec := make(map[string]Value)
			for i := 0; i < int(length); i += 2 {
				k := vm.CurrentFrame.stack[from].(String).Value
				from++
				v := vm.CurrentFrame.stack[from]
				from++
				rec[k] = v
			}
			vm.CurrentFrame.stack[to] = &Record{Value: rec}
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxErrMsg, math.MaxUint16)
		}
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
