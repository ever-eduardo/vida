package vida

import (
	"fmt"

	"github.com/alkemist-17/vida/verror"
)

type Result string

const Success Result = "Success"
const Failure Result = "Failure"

const frameSize = 1024
const stackSize = 1024

type frame struct {
	code   []uint64
	stack  []Value
	lambda *Function
	ip     int
	bp     int
	ret    int
}

type vM struct {
	Frames  [frameSize]frame
	Stack   [stackSize]Value
	Script  *Script
	Frame   *frame
	ErrInfo map[string]map[int]uint
	fp      int
}

func newVM(m *Script, extensionlibsloader LibsLoader, errInfo map[string]map[int]uint) (*vM, error) {
	extensionlibsLoader, scriptErrorInfo = extensionlibsloader, errInfo
	return &vM{Script: m, ErrInfo: errInfo}, checkISACompatibility(m)
}

func (vm *vM) run() (Result, error) {
	vm.Frame = &vm.Frames[vm.fp]
	vm.Frame.code = vm.Script.MainFunction.CoreFn.Code
	vm.Frame.lambda = vm.Script.MainFunction
	vm.Frame.stack = vm.Stack[:]
	ip := 1
	var i, op, A, B, P uint64
	for {
		i = vm.Frame.code[ip]
		op = i >> shift56
		A = i >> shift16 & clean16
		B = i & clean16
		P = i >> shift32 & clean24
		ip++
		switch op {
		case load:
			switch P {
			case loadFromLocal:
				vm.Frame.stack[B] = vm.Frame.stack[A]
			case loadFromKonst:
				vm.Frame.stack[B] = (*vm.Script.Konstants)[A]
			case loadFromGlobal:
				vm.Frame.stack[B] = (*vm.Script.Store)[A]
			default:
				vm.Frame.stack[B] = vm.Frame.lambda.Free[A]
			}
		case store:
			switch P >> shift16 {
			case storeFromGlobal:
				switch P & clean16 {
				case storeFromLocal:
					(*vm.Script.Store)[B] = vm.Frame.stack[A]
				case storeFromKonst:
					(*vm.Script.Store)[B] = (*vm.Script.Konstants)[A]
				case storeFromGlobal:
					(*vm.Script.Store)[B] = (*vm.Script.Store)[A]
				default:
					(*vm.Script.Store)[B] = vm.Frame.lambda.Free[A]
				}
			default:
				switch P & clean16 {
				case storeFromLocal:
					vm.Frame.lambda.Free[B] = vm.Frame.stack[A]
				case storeFromKonst:
					vm.Frame.lambda.Free[B] = (*vm.Script.Konstants)[A]
				case storeFromGlobal:
					vm.Frame.lambda.Free[B] = (*vm.Script.Store)[A]
				default:
					vm.Frame.lambda.Free[B] = vm.Frame.stack[A]
				}
			}
		case check:
			if P == 0 && !vm.Frame.stack[A].Boolean() {
				ip = int(B)
			}
		case jump:
			ip = int(B)
		case binopG:
			val, err := (*vm.Script.Store)[A].Binop(P>>shift16, (*vm.Script.Store)[P&clean16])
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
			val, err := vm.Frame.stack[P&clean16].Binop(P>>shift16, (*vm.Script.Konstants)[A])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case binopQ:
			val, err := (*vm.Script.Konstants)[A].Binop(P>>shift16, vm.Frame.stack[P&clean16])
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case eq:
			var val Bool
			var s byte = byte(P >> shift16)
			l := s >> shift2 & clean2bits
			r := s & clean2bits
			switch l {
			case storeFromLocal:
				switch r {
				case storeFromLocal:
					val = vm.Frame.stack[P&clean16].Equals(vm.Frame.stack[A])
				case storeFromKonst:
					val = vm.Frame.stack[P&clean16].Equals((*vm.Script.Konstants)[A])
				case storeFromGlobal:
					val = vm.Frame.stack[P&clean16].Equals((*vm.Script.Store)[A])
				default:
					val = vm.Frame.stack[P&clean16].Equals(vm.Frame.lambda.Free[A])
				}
			case storeFromKonst:
				switch r {
				case storeFromLocal:
					val = (*vm.Script.Konstants)[P&clean16].Equals(vm.Frame.stack[A])
				case storeFromGlobal:
					val = (*vm.Script.Konstants)[P&clean16].Equals((*vm.Script.Store)[A])
				default:
					val = (*vm.Script.Konstants)[P&clean16].Equals(vm.Frame.lambda.Free[A])
				}
			case storeFromGlobal:
				switch r {
				case storeFromLocal:
					val = (*vm.Script.Store)[P&clean16].Equals(vm.Frame.stack[A])
				case storeFromKonst:
					val = (*vm.Script.Store)[P&clean16].Equals((*vm.Script.Konstants)[A])
				case storeFromGlobal:
					val = (*vm.Script.Store)[P&clean16].Equals((*vm.Script.Store)[A])
				default:
					val = (*vm.Script.Store)[P&clean16].Equals(vm.Frame.lambda.Free[A])
				}
			default:
				switch r {
				case storeFromLocal:
					val = vm.Frame.lambda.Free[P&clean16].Equals(vm.Frame.stack[A])
				case storeFromKonst:
					val = vm.Frame.lambda.Free[P&clean16].Equals((*vm.Script.Konstants)[A])
				case storeFromGlobal:
					val = vm.Frame.lambda.Free[P&clean16].Equals((*vm.Script.Store)[A])
				default:
					val = vm.Frame.lambda.Free[P&clean16].Equals(vm.Frame.lambda.Free[A])
				}
			}
			if s>>4 == 1 {
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
			scopeIndex := P >> shift20
			scopeIndexable := (P >> shift16) & clean8
			switch scopeIndex {
			case storeFromLocal:
				switch scopeIndexable {
				case storeFromLocal:
					val, err = vm.Frame.stack[P&clean16].IGet(vm.Frame.stack[A])
				case storeFromGlobal:
					val, err = (*vm.Script.Store)[P&clean16].IGet(vm.Frame.stack[A])
				default:
					val, err = vm.Frame.lambda.Free[P&clean16].IGet(vm.Frame.stack[A])
				}
			case storeFromKonst:
				switch scopeIndexable {
				case storeFromLocal:
					val, err = vm.Frame.stack[P&clean16].IGet((*vm.Script.Konstants)[A])
				case storeFromGlobal:
					val, err = (*vm.Script.Store)[P&clean16].IGet((*vm.Script.Konstants)[A])
				default:
					val, err = vm.Frame.lambda.Free[P&clean16].IGet((*vm.Script.Konstants)[A])
				}
			case storeFromGlobal:
				switch scopeIndexable {
				case storeFromLocal:
					val, err = vm.Frame.stack[P&clean16].IGet((*vm.Script.Store)[A])
				case storeFromGlobal:
					val, err = (*vm.Script.Store)[P&clean16].IGet((*vm.Script.Store)[A])
				default:
					val, err = vm.Frame.lambda.Free[P&clean16].IGet((*vm.Script.Store)[A])
				}
			default:
				switch scopeIndexable {
				case storeFromLocal:
					val, err = vm.Frame.stack[P&clean16].IGet(vm.Frame.lambda.Free[A])
				case storeFromGlobal:
					val, err = (*vm.Script.Store)[P&clean16].IGet(vm.Frame.lambda.Free[A])
				default:
					val, err = vm.Frame.lambda.Free[P&clean16].IGet(vm.Frame.lambda.Free[A])
				}
			}
			if err != nil {
				return vm.createError(ip, err)
			}
			vm.Frame.stack[B] = val
		case iSet:
			var err error
			scopeIdx := P >> shift20
			scopeExp := (P >> shift16) & clean8
			switch scopeIdx {
			case storeFromLocal:
				switch scopeExp {
				case storeFromLocal:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.stack[A], vm.Frame.stack[B])
				case storeFromKonst:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.stack[A], (*vm.Script.Konstants)[B])
				case storeFromGlobal:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.stack[A], (*vm.Script.Store)[B])
				default:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.stack[A], vm.Frame.lambda.Free[B])
				}
			case storeFromKonst:
				switch scopeExp {
				case storeFromLocal:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Konstants)[A], vm.Frame.stack[B])
				case storeFromKonst:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Konstants)[A], (*vm.Script.Konstants)[B])
				case storeFromGlobal:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Konstants)[A], (*vm.Script.Store)[B])
				default:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Konstants)[A], vm.Frame.lambda.Free[B])
				}
			case storeFromGlobal:
				switch scopeExp {
				case storeFromLocal:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Store)[A], vm.Frame.stack[B])
				case storeFromKonst:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Store)[A], (*vm.Script.Konstants)[B])
				case storeFromGlobal:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Store)[A], (*vm.Script.Store)[B])
				default:
					err = vm.Frame.stack[P&clean16].ISet((*vm.Script.Store)[A], vm.Frame.lambda.Free[B])
				}
			default:
				switch scopeExp {
				case storeFromLocal:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.lambda.Free[A], vm.Frame.stack[B])
				case storeFromKonst:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.lambda.Free[A], (*vm.Script.Konstants)[B])
				case storeFromGlobal:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.lambda.Free[A], (*vm.Script.Store)[B])
				default:
					err = vm.Frame.stack[P&clean16].ISet(vm.Frame.lambda.Free[A], vm.Frame.lambda.Free[B])
				}
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
			fn := &Function{CoreFn: (*vm.Script.Konstants)[A].(*CoreFunction)}
			if fn.CoreFn.Free > 0 {
				vm.Frame.stack[B] = fn
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
					if P == ellipsisFirst {
						if xs, ok := vm.Frame.stack[B+F].(*List); ok {
							nargs = len(xs.Value) + int(F) - 1
							for i, v := range xs.Value {
								vm.Frame.stack[int(B)+int(F)+i] = v
							}
						} else {
							return vm.createError(ip, verror.ErrVariadicArgs)
						}
					} else if P == ellipsisLast {
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
			} else {
				v, err := val.Call(vm.Frame.stack[B+1 : B+A+1]...)
				if err != nil {
					return vm.createError(ip, err)
				}
				vm.Frame.stack[B] = v
			}
		case ret:
			var val Value
			switch B {
			case storeFromLocal:
				val = vm.Frame.stack[A]
			case storeFromKonst:
				val = (*vm.Script.Konstants)[A]
			case storeFromGlobal:
				val = (*vm.Script.Store)[A]
			default:
				val = vm.Frame.lambda.Free[A]
			}
			vm.fp--
			vm.Frame = &vm.Frames[vm.fp]
			ip = vm.Frame.ip
			vm.Frame.stack = vm.Stack[vm.Frame.bp:]
			vm.Frame.stack[vm.Frame.ret] = val
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("unknown opcode %v", op)
			return Failure, verror.New(vm.Frame.lambda.CoreFn.ScriptName, message, verror.RunTimeErrType, 0)
		}
	}
}

func (vm *vM) processSlice(mode, sliceable uint64) (Value, error) {
	val := vm.Frame.stack[sliceable]
	switch v := val.(type) {
	case *List:
		switch mode {
		case vcv:
			return &List{Value: v.Value[:]}, nil
		case vce:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &List{Value: v.Value[:ee]}, nil
				}
				if ee > l {
					return &List{Value: v.Value[:]}, nil
				}
				return &List{}, nil
			}
		case ecv:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &List{Value: v.Value[ee:]}, nil
				}
				if ee < 0 {
					return &List{Value: v.Value[:]}, nil
				}
				return &List{}, nil
			}
		case ece:
			l := vm.Frame.stack[sliceable+1]
			r := vm.Frame.stack[sliceable+2]
			switch ll := l.(type) {
			case Integer:
				switch rr := r.(type) {
				case Integer:
					xslen := Integer(len(v.Value))
					if ll < 0 {
						ll += xslen
					}
					if rr < 0 {
						rr += xslen
					}
					if 0 <= ll && ll <= xslen && 0 <= rr && rr <= xslen {
						return &List{Value: v.Value[ll:rr]}, nil
					}
					if ll < 0 {
						if 0 <= rr && rr <= xslen {
							return &List{Value: v.Value[:rr]}, nil
						}
						if rr > xslen {
							return &List{Value: v.Value[:]}, nil
						}
					} else if rr > xslen {
						if 0 <= ll && ll <= xslen {
							return &List{Value: v.Value[ll:]}, nil
						}
					}
				}
				return &List{}, nil
			}
		}
	case *String:
		if v.Runes == nil {
			v.Runes = []rune(v.Value)
		}
		switch mode {
		case vcv:
			return &String{Value: string(v.Runes[:])}, nil
		case vce:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &String{Value: string(v.Runes[:ee])}, nil
				}
				if ee > l {
					return &String{Value: string(v.Runes[:])}, nil
				}
				return &String{}, nil
			}
		case ecv:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &String{Value: string(v.Runes[ee:])}, nil
				}
				if ee < 0 {
					return &String{Value: string(v.Runes[:])}, nil
				}
				return &String{}, nil
			}
		case ece:
			l := vm.Frame.stack[sliceable+1]
			r := vm.Frame.stack[sliceable+2]
			switch ll := l.(type) {
			case Integer:
				switch rr := r.(type) {
				case Integer:
					xslen := Integer(len(v.Value))
					if ll < 0 {
						ll += xslen
					}
					if rr < 0 {
						rr += xslen
					}
					if 0 <= ll && ll <= xslen && 0 <= rr && rr <= xslen {
						return &String{Value: string(v.Runes[ll:rr])}, nil
					}
					if ll < 0 {
						if 0 <= rr && rr <= xslen {
							return &String{Value: string(v.Runes[:rr])}, nil
						}
						if rr > xslen {
							return &String{Value: string(v.Runes[:])}, nil
						}
					} else if rr > xslen {
						if 0 <= ll && ll <= xslen {
							return &String{Value: string(v.Runes[ll:])}, nil
						}
					}
				}
				return &String{}, nil
			}
		}
	case *Bytes:
		switch mode {
		case vcv:
			return &Bytes{Value: v.Value}, nil
		case vce:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &Bytes{Value: v.Value[:ee]}, nil
				}
				if ee > l {
					return &Bytes{Value: v.Value[:]}, nil
				}
				return &Bytes{}, nil
			}
		case ecv:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return &Bytes{Value: v.Value[ee:]}, nil
				}
				if ee < 0 {
					return &Bytes{Value: v.Value[:]}, nil
				}
				return &Bytes{}, nil
			}
		case ece:
			l := vm.Frame.stack[sliceable+1]
			r := vm.Frame.stack[sliceable+2]
			switch ll := l.(type) {
			case Integer:
				switch rr := r.(type) {
				case Integer:
					xslen := Integer(len(v.Value))
					if ll < 0 {
						ll += xslen
					}
					if rr < 0 {
						rr += xslen
					}
					if 0 <= ll && ll <= xslen && 0 <= rr && rr <= xslen {
						return &Bytes{Value: v.Value[ll:rr]}, nil
					}
					if ll < 0 {
						if 0 <= rr && rr <= xslen {
							return &Bytes{Value: v.Value[:rr]}, nil
						}
						if rr > xslen {
							return &Bytes{Value: v.Value[:]}, nil
						}
					} else if rr > xslen {
						if 0 <= ll && ll <= xslen {
							return &Bytes{Value: v.Value[ll:]}, nil
						}
					}
				}
				return &Bytes{}, nil
			}
		}
	}
	return NilValue, verror.ErrSlice
}

func (vm *vM) printCallStack() {
	fmt.Printf("  [Call Stack]\n\n")
	for i := vm.fp; i >= 0; i-- {
		modName := vm.Frames[i].lambda.CoreFn.ScriptName
		ip := vm.Frames[i].ip
		err := verror.NewStackFrameInfo(modName, vm.ErrInfo[modName][ip])
		fmt.Printf("%v\n", err)
	}
}

func (vm *vM) createError(ip int, err error) (Result, error) {
	modName := vm.Frame.lambda.CoreFn.ScriptName
	vm.Frame.ip = ip
	return Failure, verror.New(modName, err.Error(), verror.RunTimeErrType, vm.ErrInfo[modName][ip])
}

func checkISACompatibility(m *Script) error {
	majorFromCode := (m.MainFunction.CoreFn.Code[0] >> 24) & 255
	if majorFromCode == major {
		return nil
	}
	return verror.New(m.MainFunction.CoreFn.ScriptName, "script compiled with an uncompatible interpreter version", verror.FileErrType, 0)
}
