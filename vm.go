package vida

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Result string

const Success Result = "Success"
const Failure Result = "Failure"

const frameSize = 1024
const stackSize = 256
const genStackSize = 16

type frame struct {
	code  []byte
	stack []Value
	fn    *Function
	ip    int
	bp    int
	ret   byte
}

type VM struct {
	Frames  [frameSize]frame
	Stack   [stackSize]Value
	CoreLib map[string]Value
	Module  *Module
	Frame   *frame
	fp      int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m, CoreLib: loadCoreLib()}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	vm.Frame = &vm.Frames[vm.fp]
	vm.Frame.code = vm.Module.Function.Code
	vm.Frame.stack = vm.Stack[:]
	ip := 8
	for {
		op := vm.Frame.code[ip]
		ip++
		switch op {
		case setG:
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			vm.Module.Store[vm.Module.Konstants[to].(String).Value] = vm.valueFrom(scope, from)
		case setL:
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			vm.Frame.stack[to] = vm.valueFrom(scope, from)
		case move:
			from := vm.Frame.code[ip]
			ip++
			to := vm.Frame.code[ip]
			ip++
			vm.Frame.stack[to] = vm.Frame.stack[from]
		case setF:
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			vm.Frame.fn.Free[to] = vm.valueFrom(scope, from)
		case testF:
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			if !vm.valueFrom(scope, from).Boolean() {
				ip = int(jump)
			}
		case jump:
			ip = int(binary.NativeEndian.Uint16(vm.Frame.code[ip:]))
		case binop:
			op := vm.Frame.code[ip]
			ip++
			scopeLHS := vm.Frame.code[ip]
			ip++
			scopeRHS := vm.Frame.code[ip]
			ip++
			fromLHS := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromRHS := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			val, err := vm.valueFrom(scopeLHS, fromLHS).Binop(op, vm.valueFrom(scopeRHS, fromRHS))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[to] = val
		case equals:
			op := vm.Frame.code[ip]
			ip++
			scopeLHS := vm.Frame.code[ip]
			ip++
			scopeRHS := vm.Frame.code[ip]
			ip++
			fromLHS := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromRHS := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			val := vm.valueFrom(scopeLHS, fromLHS).Equals(vm.valueFrom(scopeRHS, fromRHS))
			if op == byte(token.NEQ) {
				val = !val
			}
			vm.Frame.stack[to] = val
		case prefix:
			op := vm.Frame.code[ip]
			ip++
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			val, err := vm.valueFrom(scope, from).Prefix(op)
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[to] = val
		case iGet:
			scopeIndexable := vm.Frame.code[ip]
			ip++
			scopeIndex := vm.Frame.code[ip]
			ip++
			fromIndexable := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromIndex := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			val, err := vm.valueFrom(scopeIndexable, fromIndexable).IGet(vm.valueFrom(scopeIndex, fromIndex))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[to] = val
		case iSet:
			scopeIndex := vm.Frame.code[ip]
			ip++
			scopeExpr := vm.Frame.code[ip]
			ip++
			fromIndex := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromExpr := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			from := vm.Frame.code[ip]
			ip += 2
			err := vm.valueFrom(rLoc, uint16(from)).ISet(vm.valueFrom(scopeIndex, fromIndex), vm.valueFrom(scopeExpr, fromExpr))
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
		case slice:
			mode := vm.Frame.code[ip]
			ip++
			scopeV := vm.Frame.code[ip]
			ip++
			scopeL := vm.Frame.code[ip]
			ip++
			scopeR := vm.Frame.code[ip]
			ip++
			fromV := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromL := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			fromR := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			val, err := vm.processSlice(mode, fromV, fromL, fromR, scopeV, scopeL, scopeR)
			if err != nil {
				return Failure, verror.New(vm.Module.Name, "Runtime error", verror.RunTimeErrMsg, math.MaxUint16)
			}
			vm.Frame.stack[to] = val
		case list:
			length := vm.Frame.code[ip]
			ip++
			from := vm.Frame.code[ip]
			ip++
			to := vm.Frame.code[ip]
			ip++
			xs := make([]Value, length)
			for i := 0; i < int(length); i++ {
				xs[i] = vm.Frame.stack[from]
				from++
			}
			vm.Frame.stack[to] = &List{Value: xs}
		case obj:
			length := vm.Frame.code[ip]
			ip++
			from := vm.Frame.code[ip]
			ip++
			to := vm.Frame.code[ip]
			ip++
			rec := make(map[string]Value)
			for i := 0; i < int(length); i += 2 {
				k := vm.Frame.stack[from].(String).Value
				from++
				v := vm.Frame.stack[from]
				from++
				rec[k] = v
			}
			vm.Frame.stack[to] = &Object{Value: rec}
		case forSet:
			i := vm.Frame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			if _, isInteger := vm.Frame.stack[i].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if _, isInteger := vm.Frame.stack[i+1].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if v, isInteger := vm.Frame.stack[i+2].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			} else {
				if v == 0 {
					return Failure, verror.RuntimeError
				}
			}
			ip = int(jump)
		case iForSet:
			scope := vm.Frame.code[ip]
			ip++
			reg := vm.Frame.code[ip]
			ip++
			idx := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			val := vm.valueFrom(scope, idx)
			if !val.IsIterable() {
				return Failure, verror.RuntimeError
			}
			vm.Frame.stack[reg] = val.Iterator()
			ip = int(jump)
		case forLoop:
			r := vm.Frame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			i := vm.Frame.stack[r].(Integer)
			e := vm.Frame.stack[r+1].(Integer)
			s := vm.Frame.stack[r+2].(Integer)
			if s > 0 {
				if i < e {
					vm.Frame.stack[r+3] = i
					i += s
					vm.Frame.stack[r] = i
					ip = int(jump)
				}
			} else {
				if i > e {
					vm.Frame.stack[r+3] = i
					i += s
					vm.Frame.stack[r] = i
					ip = int(jump)
				}
			}
		case iForLoop:
			r := vm.Frame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			i, _ := vm.Frame.stack[r].(Iterator)
			if i.Next() {
				vm.Frame.stack[r+1] = i.Key()
				vm.Frame.stack[r+2] = i.Value()
				ip = int(jump)
				continue
			}
		case fun:
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			to := vm.Frame.code[ip]
			ip++
			c := &Function{Core: vm.Module.Konstants[from].(*FunctionCore)}
			if c.Core.Free > 0 {
				var f []Value
				for i := 0; i < c.Core.Free; i++ {
					if c.Core.Info[i].IsLocal {
						f = append(f, vm.Frame.stack[c.Core.Info[i].Index])
					} else {
						f = append(f, vm.Frame.fn.Free[c.Core.Info[i].Index])
					}
				}
				c.Free = f
			}
			vm.Frame.stack[to] = c
		case call:
			from := vm.Frame.code[ip]
			ip++
			args := vm.Frame.code[ip]
			ip++
			val := vm.Frame.stack[from]
			if !val.IsCallable() {
				return Failure, verror.RuntimeError
			}
			if fn, ok := val.(*Function); ok {
				if args != byte(fn.Core.Arity) {
					return Failure, verror.RuntimeError
				}
				if vm.fp >= frameSize {
					return Failure, verror.RuntimeError
				}
				if fn == vm.Frame.fn && vm.Frame.code[ip] == ret {
					for i := 0; i < int(args); i++ {
						vm.Frame.stack[i] = vm.Frame.stack[int(from)+1+i]
					}
					ip = 0
					continue
				}
				vm.Frame.ip = ip
				vm.Frame.ret = from
				bs := vm.Frame.bp
				vm.fp++
				vm.Frame = &vm.Frames[vm.fp]
				vm.Frame.fn = fn
				vm.Frame.bp = bs + int(from) + 1
				vm.Frame.code = fn.Core.Code
				vm.Frame.stack = vm.Stack[vm.Frame.bp:]
				ip = 0
			} else if fn, ok := val.(GoFn); ok {
				v, err := fn(vm.Frame.stack[from+1 : from+args+1]...)
				if err != nil {
					return Failure, err
				}
				vm.Frame.stack[from] = v
			}
		case ret:
			scope := vm.Frame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.Frame.code[ip:])
			ip += 2
			val := vm.valueFrom(scope, from)
			vm.fp--
			vm.Frame = &vm.Frames[vm.fp]
			ip = vm.Frame.ip
			vm.Frame.stack = vm.Stack[vm.Frame.bp:]
			vm.Frame.stack[vm.Frame.ret] = val
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxErrMsg, math.MaxUint16)
		}
	}
}

func (vm *VM) valueFrom(scope byte, from uint16) Value {
	switch scope {
	case rKonst:
		return vm.Module.Konstants[from]
	case rLoc:
		return vm.Frame.stack[from]
	case rGlob:
		if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
	case rCore:
		if v, defined := vm.CoreLib[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
	case rFree:
		return vm.Frame.fn.Free[from]
	default:
		return NilValue
	}
}

func (vm *VM) processSlice(mode byte, fromV uint16, fromL uint16, fromR uint16, scopeV byte, scopeL byte, scopeR byte) (Value, error) {
	val := vm.valueFrom(scopeV, fromV)
	switch v := val.(type) {
	case *List:
		switch mode {
		case vcv:
			return &List{Value: v.Value[:]}, nil
		case vce:
			e := vm.valueFrom(scopeR, fromR)
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
			e := vm.valueFrom(scopeL, fromL)
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
			l := vm.valueFrom(scopeL, fromL)
			r := vm.valueFrom(scopeR, fromR)
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
	case String:
		switch mode {
		case vcv:
			return String{Value: v.Value[:]}, nil
		case vce:
			e := vm.valueFrom(scopeR, fromR)
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return String{Value: v.Value[:ee]}, nil
				}
				if ee > l {
					return String{Value: v.Value[:]}, nil
				}
				return String{}, nil
			}
		case ecv:
			e := vm.valueFrom(scopeL, fromL)
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return String{Value: v.Value[ee:]}, nil
				}
				if ee < 0 {
					return String{Value: v.Value[:]}, nil
				}
				return String{}, nil
			}
		case ece:
			l := vm.valueFrom(scopeL, fromL)
			r := vm.valueFrom(scopeR, fromR)
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
						return String{Value: v.Value[ll:rr]}, nil
					}
					if ll < 0 {
						if 0 <= rr && rr <= xslen {
							return String{Value: v.Value[:rr]}, nil
						}
						if rr > xslen {
							return String{Value: v.Value[:]}, nil
						}
					} else if rr > xslen {
						if 0 <= ll && ll <= xslen {
							return String{Value: v.Value[ll:]}, nil
						}
					}
				}
				return String{}, nil
			}
		}
	}
	return NilValue, verror.RuntimeError
}

func checkISACompatibility(m *Module) error {
	if m.Function.Code[4] == major {
		return nil
	}
	return verror.New(m.Name, "The module was compilated with another ABI version", verror.FileErrMsg, 0)
}
