package vida

import (
	"fmt"

	"github.com/alkemist-17/vida/token"
	"github.com/alkemist-17/vida/verror"
)

type ThreadState int

const (
	Ready ThreadState = iota
	Running
	Suspended
	Waiting
	Closed
)

func (state ThreadState) String() string {
	switch state {
	case Ready:
		return "ready"
	case Running:
		return "running"
	case Suspended:
		return "suspended"
	case Waiting:
		return "waiting"
	case Closed:
		return "closed"
	default:
		return "unknown"
	}
}

type Thread struct {
	ReferenceSemanticsImpl
	Frames []frame
	Stack  []Value
	Script *Script
	Frame  *frame
	State  ThreadState
	fp     int
}

func newMainThread(script *Script, extensionlibsloader LibsLoader) (*Thread, error) {
	extensionlibsLoader, clbu = extensionlibsloader, script.Store
	th := &Thread{
		Frames: make([]frame, frameSize),
		Stack:  make([]Value, fullStack),
		Script: script,
	}
	(*(script.Store))[mainThIndex] = th
	return th, nil
}

func newThread(fn *Function, script *Script, size int) *Thread {
	return &Thread{
		Script: &Script{
			Konstants:    script.Konstants,
			Store:        script.Store,
			ErrorInfo:    script.ErrorInfo,
			MainFunction: fn,
		},
		Frames: make([]frame, size),
		Stack:  make([]Value, size),
	}
}

func (th *Thread) Boolean() Bool {
	return Bool(true)
}

func (th *Thread) Prefix(op uint64) (Value, error) {
	switch op {
	case uint64(token.NOT):
		return Bool(false), nil
	default:
		return NilValue, verror.ErrPrefixOpNotDefined
	}
}

func (th *Thread) Binop(op uint64, rhs Value) (Value, error) {
	switch op {
	case uint64(token.OR):
		return th, nil
	case uint64(token.AND):
		return rhs, nil
	case uint64(token.IN):
		return IsMemberOf(th, rhs)
	}
	return NilValue, verror.ErrBinaryOpNotDefined
}

func (th *Thread) Equals(other Value) Bool {
	if val, ok := other.(*Thread); ok {
		return th == val
	}
	return false
}

func (th *Thread) String() string {
	return fmt.Sprintf("Thread(%p)", th)
}

func (th *Thread) Type() string {
	return "thread"
}

func (th *Thread) Clone() Value {
	return th
}
