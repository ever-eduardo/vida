package verror

import (
	"errors"
	"fmt"
)

const (
	FileErrType        = "file"
	LexicalErrType     = "lexical"
	SyntaxErrType      = "syntactic"
	CompilationErrType = "compilation"
	RunTimeErrType     = "runtime"
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
		return fmt.Sprintf("\n\n[%v]\n   Module  : %v\n   Line    : %v\n   Message : %v\n\n", e.ErrType, e.ModuleName, e.Line, e.Message)
	default:
		if e.Line == 0 {
			return fmt.Sprintf("\n\n[Error]\n   Type    : %v\n   Module  : %v\n   Message : %v\n\n", e.ErrType, e.ModuleName, e.Message)
		}
		return fmt.Sprintf("\n\n[Error]\n   Type    : %v\n   Module  : %v\n   Line    : %v\n   Message : %v\n\n", e.ErrType, e.ModuleName, e.Line, e.Message)
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
	ErrRuntime     = errors.New(RunTimeErrType)
	ErrStringLimit = errors.New("strings max size has been reached")
	ErrCompilation = errors.New(CompilationErrType)
)
