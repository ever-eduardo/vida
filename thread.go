package vida

import (
	"fmt"

	"github.com/alkemist-17/vida/token"
	"github.com/alkemist-17/vida/verror"
)

type Thread struct {
	ReferenceSemanticsImpl
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
