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

var (
	ErrStringLimit              = errors.New("strings max size has been reached")
	ErrOpNotDefinedForIterators = errors.New("operation not defined for iterators")
	ErrValueNotIndexable        = errors.New("trying to index a non indexable value")
	ErrPrefixOpNotDefined       = errors.New("prefix operation not defined")
	ErrBinaryOpNotDefined       = errors.New("binary operation not defined")
	ErrDivisionByZero           = errors.New("division by zero not defined")
)
