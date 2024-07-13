package vida

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

func (vm *VM) Inspect(ip int) {
	clear()
	fmt.Println("Running", vm.Module.Name)
	fmt.Printf("Store => %v\n", vm.Module.Store)
	fmt.Print("Konst => ")
	for i, v := range vm.Module.Konstants {
		fmt.Printf("[%v: %v], ", i, v)
	}
	fmt.Println()
	fmt.Printf("Instr => %v\n", printInstr(ip, vm.Module.Code))
	fmt.Println("Stack =>")
	for i, v := range vm.Stack {
		if v != nil {
			fmt.Printf("  [%v] %v\n", i, v)
		}
	}
	fmt.Printf("\nPress 'Enter' to continue => ")
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
		case equals:
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
			val := vm.valueFrom(scopeLHS, fromLHS).Equals(vm.valueFrom(scopeRHS, fromRHS))
			if op == byte(token.NEQ) {
				val = !val
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
		case iSet:
			scopeIndex := vm.CurrentFrame.code[ip]
			ip++
			scopeExpr := vm.CurrentFrame.code[ip]
			ip++
			fromIndex := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			fromExpr := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			from := vm.CurrentFrame.code[ip]
			ip += 2
			err := vm.valueFrom(rLocal, uint16(from)).ISet(vm.valueFrom(scopeIndex, fromIndex), vm.valueFrom(scopeExpr, fromExpr))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
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
		case document:
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
			vm.CurrentFrame.stack[to] = &Document{Value: rec}
		case forInit:
			forIdx := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			forLoop := vm.valueFrom(rKonst, forIdx).(*ForLoopState)
			if _, isInteger := vm.CurrentFrame.stack[forLoop.Init].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if _, isInteger := vm.CurrentFrame.stack[forLoop.End].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if v, isInteger := vm.CurrentFrame.stack[forLoop.Step].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			} else {
				if v == 0 {
					return Failure, verror.RuntimeError
				}
			}
			ip = int(jump)
		case forLoop:
			forIdx := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			forLoop := vm.valueFrom(rKonst, forIdx).(*ForLoopState)
			i := vm.CurrentFrame.stack[forLoop.Init].(Integer)
			e := vm.CurrentFrame.stack[forLoop.End].(Integer)
			s := vm.CurrentFrame.stack[forLoop.Step].(Integer)
			if s > 0 {
				if i < e {
					vm.CurrentFrame.stack[forLoop.State] = i
					i += s
					vm.CurrentFrame.stack[forLoop.Init] = i
					ip = int(jump)
				}
			} else {
				if i > e {
					vm.CurrentFrame.stack[forLoop.State] = i
					i += s
					vm.CurrentFrame.stack[forLoop.Init] = i
					ip = int(jump)
				}
			}
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
