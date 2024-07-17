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

const callStackSize = 1024
const stackSize = 256

type frame struct {
	code  []byte
	stack []Value
	ip    int
	ret   int
	op    byte
}

type VM struct {
	Frames       [callStackSize]frame
	Stack        [stackSize]Value
	Prelude      map[string]Value
	Module       *Module
	CurrentFrame *frame
	fp           int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m, Prelude: loadPrelude()}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	vm.CurrentFrame = &vm.Frames[vm.fp]
	vm.CurrentFrame.code = vm.Module.Code
	vm.CurrentFrame.stack = vm.Stack[:]
	ip := 8
	for {
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
		case testF:
			scope := vm.CurrentFrame.code[ip]
			ip++
			from := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			if !vm.valueFrom(scope, from).Boolean() {
				ip = int(jump)
			}
		case jump:
			ip = int(binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:]))
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
		case doc:
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
		case forSet:
			i := vm.CurrentFrame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			if _, isInteger := vm.CurrentFrame.stack[i].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if _, isInteger := vm.CurrentFrame.stack[i+1].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			}
			if v, isInteger := vm.CurrentFrame.stack[i+2].(Integer); !isInteger {
				return Failure, verror.RuntimeError
			} else {
				if v == 0 {
					return Failure, verror.RuntimeError
				}
			}
			ip = int(jump)
		case iForSet:
			scope := vm.CurrentFrame.code[ip]
			ip++
			reg := vm.CurrentFrame.code[ip]
			ip++
			idx := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			val := vm.valueFrom(scope, idx)
			if !val.IsIterable() {
				return Failure, verror.RuntimeError
			}
			vm.CurrentFrame.stack[reg] = val.Iterator()
			ip = int(jump)
		case forLoop:
			r := vm.CurrentFrame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			i := vm.CurrentFrame.stack[r].(Integer)
			e := vm.CurrentFrame.stack[r+1].(Integer)
			s := vm.CurrentFrame.stack[r+2].(Integer)
			if s > 0 {
				if i < e {
					vm.CurrentFrame.stack[r+3] = i
					i += s
					vm.CurrentFrame.stack[r] = i
					ip = int(jump)
				}
			} else {
				if i > e {
					vm.CurrentFrame.stack[r+3] = i
					i += s
					vm.CurrentFrame.stack[r] = i
					ip = int(jump)
				}
			}
		case iForLoop:
			r := vm.CurrentFrame.code[ip]
			ip++
			jump := binary.NativeEndian.Uint16(vm.CurrentFrame.code[ip:])
			ip += 2
			i, _ := vm.CurrentFrame.stack[r].(Iterator)
			if i.Next() {
				vm.CurrentFrame.stack[r+1] = i.Key()
				vm.CurrentFrame.stack[r+2] = i.Value()
				ip = int(jump)
				continue
			}
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
	case rLocal:
		return vm.CurrentFrame.stack[from]
	case rGlobal:
		if v, defined := vm.Module.Store[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
	case rCore:
		if v, defined := vm.Prelude[vm.Module.Konstants[from].(String).Value]; defined {
			return v
		} else {
			return NilValue
		}
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
	if m.Code[4] == major {
		return nil
	}
	return verror.New(m.Name, "The module was compilated with another ABI version", verror.FileErrMsg, 0)
}
