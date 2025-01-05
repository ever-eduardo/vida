package verror

import (
	"errors"
	"fmt"
)

const (
	FileErrType        = "File"
	LexicalErrType     = "Lexical"
	SyntaxErrType      = "Syntax"
	CompilationErrType = "Compilation"
	RunTimeErrType     = "Runtime"
	AssertionErrType   = "Assertion Failure"
	ExceptionErrType   = "Exception"
)

type VidaError struct {
	ModuleName string
	Message    string
	ErrType    string
	Line       uint
}

func (e VidaError) Error() string {
	switch e.ErrType {
	case ExceptionErrType, AssertionErrType:
		return fmt.Sprintf("\n\n  [%v]\n   Module    : %v\n   Near line : %v\n   Message   : %v\n\n", e.ErrType, e.ModuleName, e.Line, e.Message)
	default:
		if e.Line == 0 {
			return fmt.Sprintf("\n\n  [%v Error]\n   Module  : %v\n   Message : %v\n\n", e.ErrType, e.ModuleName, e.Message)
		}
		return fmt.Sprintf("\n\n  [%v Error]\n   Module    : %v\n   Near line : %v\n   Message   : %v\n\n", e.ErrType, e.ModuleName, e.Line, e.Message)
	}
}

func New(moduleName string, message string, errorType string, line uint) VidaError {
	return VidaError{
		ModuleName: moduleName,
		Line:       line,
		Message:    message,
		ErrType:    errorType,
	}
}

type StackFrameInfo struct {
	ModuleName string
	Line       uint
}

func NewStackFrameInfo(moduleName string, line uint) StackFrameInfo {
	return StackFrameInfo{
		ModuleName: moduleName,
		Line:       line,
	}
}

func (sfi StackFrameInfo) Error() string {
	return fmt.Sprintf("   Module    : %v\n   Near line : %v\n", sfi.ModuleName, sfi.Line)
}

var (
	ErrStringLimit                      = errors.New("strings max size has been reached")
	ErrOpNotDefinedForIterators         = errors.New("operation not defined for iterators")
	ErrValueNotIndexable                = errors.New("value is not indexable")
	ErrPrefixOpNotDefined               = errors.New("prefix operation not defined")
	ErrBinaryOpNotDefined               = errors.New("binary operation not defined")
	ErrDivisionByZero                   = errors.New("division by zero not defined")
	ErrExpectedInteger                  = errors.New("expected a value of type integer")
	ErrExpectedIntegerDifferentFromZero = errors.New("expected an integer value different from zero")
	ErrValueNotIterable                 = errors.New("value is not iterable")
	ErrValueNotCallable                 = errors.New("value is not callable")
	ErrStackOverflow                    = errors.New("stack overflow")
	ErrArity                            = errors.New("given arguments count is different from arity definition")
	ErrNotEnoughArgs                    = errors.New("not given enough arguments to the function")
	ErrVariadicArgs                     = errors.New("expected a list for variradic arguments")
	ErrSlice                            = errors.New("could not process the slice")
	ErrValueIsConstant                  = errors.New("value is constant")
)
