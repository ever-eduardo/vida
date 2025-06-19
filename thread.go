package vida

import (
	"fmt"

	"github.com/alkemist-17/vida/token"
	"github.com/alkemist-17/vida/verror"
)

type ThreadState int

const (
	Created ThreadState = iota
	Running
	Suspended
	Closed
)

func (state ThreadState) String() string {
	switch state {
	case Created:
		return "created"
	case Running:
		return "running"
	case Suspended:
		return "suspended"
	case Closed:
		return "closed"
	default:
		return "unknown"
	}
}

type Suspendable interface {
	Resume(args ...Value) Value
	State() ThreadState
	IsAlive() bool
	Suspend(args ...Value) Value
}

type Thread struct {
	ReferenceSemanticsImpl
	Frames  []frame
	Stack   []Value
	Script  *Script
	Frame   *frame
	ErrInfo map[string]map[int]uint
	State   ThreadState
	fp      int
}

func createThread(fn *Function, store *[]Value, konstants *[]Value, size int) *Thread {
	script := &Script{
		Konstants:    konstants,
		Store:        store,
		MainFunction: fn,
	}
	th := &Thread{
		Script:  script,
		ErrInfo: scriptErrorInfo,
		Frames:  make([]frame, size),
		Stack:   make([]Value, size),
	}
	return th
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
