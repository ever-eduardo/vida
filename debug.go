package vida

import (
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/verror"
)

func (vm *VM) Inspect(ip int) {
	clear()
	fmt.Println("Running", vm.Module.Name)
	fmt.Printf("Store => ")
	for i := len(coreLibNames); i < len(vm.Module.Store); i++ {
		fmt.Printf("[%v -> %v], ", i, vm.Module.Store[i])
	}
	fmt.Println()
	fmt.Print("Konst => ")
	for i, v := range vm.Module.Konstants {
		fmt.Printf("[%v -> %v], ", i, v)
	}
	fmt.Println()
	fmt.Printf("Frame => %v\n", vm.fp)
	fmt.Printf("Ip    => %v\n", ip)
	s := printInstr(vm.Frame.code[ip], uint64(ip), true)
	fmt.Printf("Instr => %v\n", s)
	fmt.Println("Stack =>")
	for i, v := range vm.Stack {
		if v != nil {
			if vm.Frame.bp == i {
				fmt.Printf(" *[%v] %v\n", i, v)
			} else {
				fmt.Printf("  [%v] %v\n", i, v)
			}
		}
	}
	fmt.Printf("\nPress 'Enter' to continue => ")
	fmt.Scanf(" ")
}

func (vm *VM) Debug() (Result, error) {
	vm.Frame = &vm.Frames[vm.fp]
	vm.Frame.code = vm.Module.MainFunction.CoreFn.Code
	vm.Frame.lambda = vm.Module.MainFunction
	vm.Frame.stack = vm.Stack[:]
	ip := 1
	var i, op, A, B, P uint64
	for {
		vm.Inspect(ip)
		i = vm.Frame.code[ip]
		op = i >> shift56
		A = i >> shift16 & clean16
		B = i & clean16
		P = i >> shift32 & clean24
		ip++
		switch op {
		case storeG:
			if P == 1 {
				vm.Module.Store[B] = vm.Module.Konstants[A]
			} else {
				vm.Module.Store[B] = vm.Frame.stack[A]
			}
		case loadG:
			vm.Frame.stack[B] = vm.Module.Store[A]
		case loadF:
			vm.Frame.stack[B] = vm.Frame.lambda.Free[A]
		case loadK:
			vm.Frame.stack[B] = vm.Module.Konstants[A]
		case move:
			vm.Frame.stack[B] = vm.Frame.stack[A]
		case storeF:
			vm.Frame.lambda.Free[B] = vm.Frame.stack[A]
		// case checkF:
		// 	scope := vm.Frame.code[ip]
		// 	ip++
		// 	from := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	jump := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	if !vm.valueFrom(scope, from).Boolean() {
		// 		ip = int(jump)
		// 	}
		// case jump:
		// 	ip = int(uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8)
		case binopG:
			val, err := vm.Module.Store[A].Binop(P>>shift16, vm.Module.Store[P&clean16])
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[B] = val
		case binop:
			val, err := vm.Frame.stack[A].Binop(P>>shift16, vm.Frame.stack[P&clean16])
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[B] = val
		case binopK:
			val, err := vm.Module.Konstants[A].Binop(P>>shift16, vm.Frame.stack[P&clean16])
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[B] = val
		// case equals:
		// 	op := vm.Frame.code[ip]
		// 	ip++
		// 	scopeLHS := vm.Frame.code[ip]
		// 	ip++
		// 	scopeRHS := vm.Frame.code[ip]
		// 	ip++
		// 	fromLHS := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	fromRHS := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	val := vm.valueFrom(scopeLHS, fromLHS).Equals(vm.valueFrom(scopeRHS, fromRHS))
		// 	if op == byte(token.NEQ) {
		// 		val = !val
		// 	}
		// 	vm.Frame.stack[to] = val
		case prefix:
			val, err := vm.Frame.stack[A].Prefix(P)
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[B] = val
		// case iGet:
		// 	scopeIndexable := vm.Frame.code[ip]
		// 	ip++
		// 	scopeIndex := vm.Frame.code[ip]
		// 	ip++
		// 	fromIndexable := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	fromIndex := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	val, err := vm.valueFrom(scopeIndexable, fromIndexable).IGet(vm.valueFrom(scopeIndex, fromIndex))
		// 	if err != nil {
		// 		return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
		// 	}
		// 	vm.Frame.stack[to] = val
		// case iSet:
		// 	scopeIndex := vm.Frame.code[ip]
		// 	ip++
		// 	scopeExpr := vm.Frame.code[ip]
		// 	ip++
		// 	fromIndex := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	fromExpr := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	from := vm.Frame.code[ip]
		// 	ip += 2
		// 	err := vm.valueFrom(rLoc, uint16(from)).ISet(vm.valueFrom(scopeIndex, fromIndex), vm.valueFrom(scopeExpr, fromExpr))
		// 	if err != nil {
		// 		return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
		// 	}
		// case slice:
		// 	mode := vm.Frame.code[ip]
		// 	ip++
		// 	scopeV := vm.Frame.code[ip]
		// 	ip++
		// 	scopeL := vm.Frame.code[ip]
		// 	ip++
		// 	scopeR := vm.Frame.code[ip]
		// 	ip++
		// 	fromV := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	fromL := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	fromR := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	val, err := vm.processSlice(mode, fromV, fromL, fromR, scopeV, scopeL, scopeR)
		// 	if err != nil {
		// 		return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
		// 	}
		// 	vm.Frame.stack[to] = val
		// case list:
		// 	length := vm.Frame.code[ip]
		// 	ip++
		// 	from := vm.Frame.code[ip]
		// 	ip++
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	xs := make([]Value, length)
		// 	for i := 0; i < int(length); i++ {
		// 		xs[i] = vm.Frame.stack[from]
		// 		from++
		// 	}
		// 	vm.Frame.stack[to] = &List{Value: xs}
		// case obj:
		// 	length := vm.Frame.code[ip]
		// 	ip++
		// 	from := vm.Frame.code[ip]
		// 	ip++
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	rec := make(map[string]Value)
		// 	for i := 0; i < int(length); i += 2 {
		// 		k := vm.Frame.stack[from].(String).Value
		// 		from++
		// 		v := vm.Frame.stack[from]
		// 		from++
		// 		rec[k] = v
		// 	}
		// 	vm.Frame.stack[to] = &Object{Value: rec}
		// case forSet:
		// 	i := vm.Frame.code[ip]
		// 	ip++
		// 	jump := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	if _, isInteger := vm.Frame.stack[i].(Integer); !isInteger {
		// 		return Failure, verror.RuntimeError
		// 	}
		// 	if _, isInteger := vm.Frame.stack[i+1].(Integer); !isInteger {
		// 		return Failure, verror.RuntimeError
		// 	}
		// 	if v, isInteger := vm.Frame.stack[i+2].(Integer); !isInteger {
		// 		return Failure, verror.RuntimeError
		// 	} else {
		// 		if v == 0 {
		// 			return Failure, verror.RuntimeError
		// 		}
		// 	}
		// 	ip = int(jump)
		// case iForSet:
		// 	scope := vm.Frame.code[ip]
		// 	ip++
		// 	reg := vm.Frame.code[ip]
		// 	ip++
		// 	idx := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	jump := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	val := vm.valueFrom(scope, idx)
		// 	if !val.IsIterable() {
		// 		return Failure, verror.RuntimeError
		// 	}
		// 	vm.Frame.stack[reg] = val.Iterator()
		// 	ip = int(jump)
		// case forLoop:
		// 	r := vm.Frame.code[ip]
		// 	ip++
		// 	jump := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	i := vm.Frame.stack[r].(Integer)
		// 	e := vm.Frame.stack[r+1].(Integer)
		// 	s := vm.Frame.stack[r+2].(Integer)
		// 	if s > 0 {
		// 		if i < e {
		// 			vm.Frame.stack[r+3] = i
		// 			i += s
		// 			vm.Frame.stack[r] = i
		// 			ip = int(jump)
		// 		}
		// 	} else {
		// 		if i > e {
		// 			vm.Frame.stack[r+3] = i
		// 			i += s
		// 			vm.Frame.stack[r] = i
		// 			ip = int(jump)
		// 		}
		// 	}
		// case iForLoop:
		// 	r := vm.Frame.code[ip]
		// 	ip++
		// 	jump := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	i, _ := vm.Frame.stack[r].(Iterator)
		// 	if i.Next() {
		// 		vm.Frame.stack[r+1] = i.Key()
		// 		vm.Frame.stack[r+2] = i.Value()
		// 		ip = int(jump)
		// 		continue
		// 	}
		// case fun:
		// 	from := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	to := vm.Frame.code[ip]
		// 	ip++
		// 	fn := &Function{CoreFn: vm.Module.Konstants[from].(*CoreFunction)}
		// 	if fn.CoreFn.Free > 0 {
		// 		var free []Value
		// 		for i := 0; i < fn.CoreFn.Free; i++ {
		// 			if fn.CoreFn.Info[i].IsLocal {
		// 				free = append(free, vm.Frame.stack[fn.CoreFn.Info[i].Index])
		// 			} else {
		// 				free = append(free, vm.Frame.lambda.Free[fn.CoreFn.Info[i].Index])
		// 			}
		// 		}
		// 		fn.Free = free
		// 	}
		// 	vm.Frame.stack[to] = fn
		// case call:
		// 	from := vm.Frame.code[ip]
		// 	ip++
		// 	args := vm.Frame.code[ip]
		// 	ip++
		// 	val := vm.Frame.stack[from]
		// 	if !val.IsCallable() {
		// 		return Failure, verror.RuntimeError
		// 	}
		// 	if fn, ok := val.(*Function); ok {
		// 		if args != byte(fn.CoreFn.Arity) {
		// 			return Failure, verror.RuntimeError
		// 		}
		// 		if vm.fp >= frameSize {
		// 			return Failure, verror.RuntimeError
		// 		}
		// 		if fn == vm.Frame.lambda && vm.Frame.code[ip] == ret {
		// 			for i := 0; i < int(args); i++ {
		// 				vm.Frame.stack[i] = vm.Frame.stack[int(from)+1+i]
		// 			}
		// 			ip = 0
		// 			continue
		// 		}
		// 		vm.Frame.ip = ip
		// 		vm.Frame.ret = from
		// 		bs := vm.Frame.bp
		// 		vm.fp++
		// 		vm.Frame = &vm.Frames[vm.fp]
		// 		vm.Frame.lambda = fn
		// 		vm.Frame.bp = bs + int(from) + 1
		// 		vm.Frame.code = fn.CoreFn.Code
		// 		vm.Frame.stack = vm.Stack[vm.Frame.bp:]
		// 		ip = 0
		// 	} else if fn, ok := val.(GFn); ok {
		// 		v, err := fn(vm.Frame.stack[from+1 : from+args+1]...)
		// 		if err != nil {
		// 			return Failure, err
		// 		}
		// 		vm.Frame.stack[from] = v
		// 	}
		// case ret:
		// 	scope := vm.Frame.code[ip]
		// 	ip++
		// 	from := uint16(vm.Frame.code[ip]) | uint16(vm.Frame.code[ip+1])<<8
		// 	ip += 2
		// 	val := vm.valueFrom(scope, from)
		// 	vm.fp--
		// 	vm.Frame = &vm.Frames[vm.fp]
		// 	ip = vm.Frame.ip
		// 	vm.Frame.stack = vm.Stack[vm.Frame.bp:]
		// 	vm.Frame.stack[vm.Frame.ret] = val
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
