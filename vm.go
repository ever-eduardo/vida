package vida

import (
	"fmt"
	"math"

	"github.com/ever-eduardo/vida/verror"
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

type VM struct {
	Frames [frameSize]frame
	Stack  [stackSize]Value
	Module *Module
	Frame  *frame
	fp     int
}

func NewVM(m *Module) (*VM, error) {
	return &VM{Module: m}, checkISACompatibility(m)
}

func (vm *VM) Run() (Result, error) {
	vm.Frame = &vm.Frames[vm.fp]
	vm.Frame.code = vm.Module.MainFunction.CoreFn.Code
	vm.Frame.lambda = vm.Module.MainFunction
	vm.Frame.stack = vm.Stack[:]
	ip := 8
	for {
		op := vm.Frame.code[ip]
		ip++
		switch op {
		case end:
			return Success, nil
		default:
			message := fmt.Sprintf("Unknown vm instruction %v", ip)
			return Failure, verror.New(vm.Module.Name, message, verror.SyntaxErrMsg, math.MaxUint16)
		}
	}
}

func (vm *VM) processSlice(mode, sliceable uint64) (Value, error) {
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
	case String:
		if v.Runes == nil {
			v.Runes = []rune(v.Value)
		}
		switch mode {
		case vcv:
			return String{Value: string(v.Runes[:])}, nil
		case vce:
			e := vm.Frame.stack[sliceable+1]
			switch ee := e.(type) {
			case Integer:
				l := Integer(len(v.Value))
				if ee < 0 {
					ee += l
				}
				if 0 <= ee && ee <= l {
					return String{Value: string(v.Runes[:ee])}, nil
				}
				if ee > l {
					return String{Value: string(v.Runes[:])}, nil
				}
				return String{}, nil
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
					return String{Value: string(v.Runes[ee:])}, nil
				}
				if ee < 0 {
					return String{Value: string(v.Runes[:])}, nil
				}
				return String{}, nil
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
						return String{Value: string(v.Runes[ll:rr])}, nil
					}
					if ll < 0 {
						if 0 <= rr && rr <= xslen {
							return String{Value: string(v.Runes[:rr])}, nil
						}
						if rr > xslen {
							return String{Value: string(v.Runes[:])}, nil
						}
					} else if rr > xslen {
						if 0 <= ll && ll <= xslen {
							return String{Value: string(v.Runes[ll:])}, nil
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
	majorFromCode := (m.MainFunction.CoreFn.Code[0] >> 24) & 255
	if majorFromCode == major {
		return nil
	}
	return verror.New(m.Name, "The module was compiled with another ABI version", verror.FileErrMsg, 0)
}
