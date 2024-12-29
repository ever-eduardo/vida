package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

func (vm *vM) Inspect(ip int) {
	clear()
	fmt.Println("Running", vm.Frame.lambda.CoreFn.ModuleName)
	fmt.Printf("Store => ")
	for i := len(coreLibNames); i < len((*vm.Module.Store)); i++ {
		fmt.Printf("[%v -> %v], ", i, (*vm.Module.Store)[i])
	}
	fmt.Println()
	fmt.Print("Konst => ")
	for i, v := range *vm.Module.Konstants {
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

func (vm *vM) debug() (Result, error) {
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
				(*vm.Module.Store)[B] = (*vm.Module.Konstants)[A]
			} else {
				(*vm.Module.Store)[B] = vm.Frame.stack[A]
			}
		case loadG:
			vm.Frame.stack[B] = (*vm.Module.Store)[A]
		case loadF:
			vm.Frame.stack[B] = vm.Frame.lambda.Free[A]
		case loadK:
			vm.Frame.stack[B] = (*vm.Module.Konstants)[A]
		case move:
			vm.Frame.stack[B] = vm.Frame.stack[A]
		case storeF:
			vm.Frame.lambda.Free[B] = vm.Frame.stack[A]
		case check:
			if P == 0 && !vm.Frame.stack[A].Boolean() {
				ip = int(B)
			}
		case jump:
			ip = int(B)
		case binopG:
			val, err := (*vm.Module.Store)[A].Binop(P>>shift16, (*vm.Module.Store)[P&clean16])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case binop:
			val, err := vm.Frame.stack[A].Binop(P>>shift16, vm.Frame.stack[P&clean16])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case binopK:
			val, err := vm.Frame.stack[P&clean16].Binop(P>>shift16, (*vm.Module.Konstants)[A])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case binopQ:
			val, err := (*vm.Module.Konstants)[A].Binop(P>>shift16, vm.Frame.stack[P&clean16])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case eq:
			val := vm.Frame.stack[A].Equals(vm.Frame.stack[P&clean16])
			if P>>shift16 == uint64(token.NEQ) {
				val = !val
			}
			vm.Frame.stack[B] = val
		case eqG:
			val := (*vm.Module.Store)[A].Equals((*vm.Module.Store)[P&clean16])
			if P>>shift16 == uint64(token.NEQ) {
				val = !val
			}
			vm.Frame.stack[B] = val
		case eqK:
			val := vm.Frame.stack[P&clean16].Equals((*vm.Module.Konstants)[A])
			if P>>shift16 == uint64(token.NEQ) {
				val = !val
			}
			vm.Frame.stack[B] = val
		case eqQ:
			val := (*vm.Module.Konstants)[A].Equals(vm.Frame.stack[P&clean16])
			if P>>shift16 == uint64(token.NEQ) {
				val = !val
			}
			vm.Frame.stack[B] = val
		case prefix:
			val, err := vm.Frame.stack[A].Prefix(P)
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case iGet:
			var val Value
			var err error
			if P>>shift16 == 0 {
				val, err = vm.Frame.stack[P].IGet(vm.Frame.stack[A])
			} else {
				val, err = vm.Frame.stack[P&clean16].IGet((*vm.Module.Konstants)[A])
			}
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case iSet:
			var err error
			if P>>shift16 == 0 {
				err = vm.Frame.stack[P].ISet(vm.Frame.stack[A], vm.Frame.stack[B])
			} else {
				err = vm.Frame.stack[P&clean16].ISet(vm.Frame.stack[A], (*vm.Module.Konstants)[B])
			}
			if err != nil {
				return vm.createError(ip, err)
			}
		case iSetK:
			var err error
			if P>>shift16 == 0 {
				err = vm.Frame.stack[P].ISet((*vm.Module.Konstants)[A], vm.Frame.stack[B])
			} else {
				err = vm.Frame.stack[P&clean16].ISet((*vm.Module.Konstants)[A], (*vm.Module.Konstants)[B])
			}
			if err != nil {
				return vm.createError(ip, err)
			}
		case slice:
			val, err := vm.processSlice(P, A)
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case list:
			xs := make([]Value, P)
			F := A
			for i := 0; i < int(P); i++ {
				xs[i] = vm.Frame.stack[F]
				F++
			}
			vm.Frame.stack[B] = &List{Value: xs}
		case object:
			vm.Frame.stack[B] = &Object{Value: make(map[string]Value)}
		case forSet:
			if _, isInteger := vm.Frame.stack[B].(Integer); !isInteger {
				return vm.createError(ip, verror.ErrExpectedInteger)
			}
			if _, isInteger := vm.Frame.stack[B+1].(Integer); !isInteger {
				return vm.createError(ip, verror.ErrExpectedInteger)
			}
			if v, isInteger := vm.Frame.stack[B+2].(Integer); !isInteger {
				return vm.createError(ip, verror.ErrExpectedInteger)
			} else if v == 0 {
				return vm.createError(ip, verror.ErrExpectedIntegerDifferentFromZero)
			}
			ip = int(A)
		case iForSet:
			iterable := vm.Frame.stack[A]
			if !iterable.IsIterable() {
				return vm.createError(ip, verror.ErrValueNotIterable)
			}
			vm.Frame.stack[B] = iterable.Iterator()
			ip = int(P)
		case forLoop:
			i := vm.Frame.stack[B].(Integer)
			e := vm.Frame.stack[B+1].(Integer)
			s := vm.Frame.stack[B+2].(Integer)
			if s > 0 {
				if i < e {
					vm.Frame.stack[B+3] = i
					i += s
					vm.Frame.stack[B] = i
					ip = int(A)
				}
			} else {
				if i > e {
					vm.Frame.stack[B+3] = i
					i += s
					vm.Frame.stack[B] = i
					ip = int(A)
				}
			}
		case iForLoop:
			i, _ := vm.Frame.stack[B].(Iterator)
			if i.Next() {
				vm.Frame.stack[B+1] = i.Key()
				vm.Frame.stack[B+2] = i.Value()
				ip = int(A)
				continue
			}
		case fun:
			fn := &Function{CoreFn: (*vm.Module.Konstants)[A].(*CoreFunction)}
			if fn.CoreFn.Free > 0 {
				var free []Value
				for i := 0; i < fn.CoreFn.Free; i++ {
					if fn.CoreFn.Info[i].IsLocal {
						free = append(free, vm.Frame.stack[fn.CoreFn.Info[i].Index])
					} else {
						free = append(free, vm.Frame.lambda.Free[fn.CoreFn.Info[i].Index])
					}
				}
				fn.Free = free
			}
			vm.Frame.stack[B] = fn
		case call:
			val := vm.Frame.stack[B]
			nargs := int(A)
			F := P >> shift16
			P = P & clean16
			if !val.IsCallable() {
				return vm.createError(ip, verror.ErrValueNotCallable)
			}
			if fn, ok := val.(*Function); ok {
				if vm.fp >= frameSize {
					return vm.createError(ip, verror.ErrStackOverflow)
				}
				if P != 0 {
					if P == 1 {
						if xs, ok := vm.Frame.stack[B+F].(*List); ok {
							nargs = len(xs.Value) + int(F) - 1
							for i, v := range xs.Value {
								vm.Frame.stack[int(B)+int(F)+i] = v
							}
						} else {
							return vm.createError(ip, verror.ErrVariadicArgs)
						}
					} else if P == 2 {
						if xs, ok := vm.Frame.stack[int(B)+nargs].(*List); ok {
							nargs += len(xs.Value) - 1
							for i, v := range xs.Value {
								vm.Frame.stack[int(B)+int(A)+i] = v
							}
						} else {
							return vm.createError(ip, verror.ErrVariadicArgs)
						}
					}
				}
				if fn.CoreFn.IsVar {
					if fn.CoreFn.Arity > nargs {
						return vm.createError(ip, verror.ErrNotEnoughArgs)
					}
					init := int(B) + 1 + fn.CoreFn.Arity
					count := nargs - fn.CoreFn.Arity
					xs := make([]Value, count)
					for i := 0; i < count; i++ {
						xs[i] = vm.Frame.stack[init+i]
					}
					vm.Frame.stack[init] = &List{Value: xs}
				} else if nargs != fn.CoreFn.Arity {
					return vm.createError(ip, verror.ErrArity)
				}
				if fn == vm.Frame.lambda && vm.Frame.code[ip]>>shift56 == ret {
					for i := 0; i < nargs; i++ {
						vm.Frame.stack[i] = vm.Frame.stack[int(B)+1+i]
					}
					ip = 0
					continue
				}
				vm.Frame.ip = ip
				vm.Frame.ret = int(B)
				bs := vm.Frame.bp
				vm.fp++
				vm.Frame = &vm.Frames[vm.fp]
				vm.Frame.lambda = fn
				vm.Frame.bp = bs + int(B) + 1
				vm.Frame.code = fn.CoreFn.Code
				vm.Frame.stack = vm.Stack[vm.Frame.bp:]
				ip = 0
			} else if fn, ok := val.(GFn); ok {
				v, err := fn(vm.Frame.stack[B+1 : B+A+1]...)
				if err != nil {
					return vm.createError(ip, err)
				}
				vm.Frame.stack[B] = v
			}
		case ret:
			val := vm.Frame.stack[B]
			vm.fp--
			vm.Frame = &vm.Frames[vm.fp]
			ip = vm.Frame.ip
			vm.Frame.stack = vm.Stack[vm.Frame.bp:]
			vm.Frame.stack[vm.Frame.ret] = val
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("unknown opcode %v", op)
			return Failure, verror.New(vm.Frame.lambda.CoreFn.ModuleName, message, verror.RunTimeErrType, 0)
		}
	}
}

func clear() {
	fmt.Printf("\u001B[H")
	fmt.Printf("\u001B[2J")
}
