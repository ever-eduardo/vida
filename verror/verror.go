package verror

import (
	"errors"
	"fmt"
)

const FileErrMsg = "file"
const LexicalErrMsg = "lexical"
const SyntaxErrMsg = "syntax"
const CompilerErrorMsg = "compiler error"
const RunTimeErrMsg = "runtime error"
const AssertionErr = "assertion failure"
const FatalFailure = "fatal failure"

type VidaError struct {
	ModuleName   string
	Message      string
	Type         string
	Line         uint
	FromCompiler bool
}

func (e VidaError) Error() string {
	if e.FromCompiler {
		return fmt.Sprintf("[%v Error] : %v", e.Type, e.Message)
	}
	return fmt.Sprintf("[%v Error] : File '%v' line %v : %v", e.Type, e.ModuleName, e.Line, e.Message)
}

func New(moduleName string, message string, errorType string, line uint) VidaError {
	return VidaError{
		ModuleName:   moduleName,
		Line:         line,
		Message:      message,
		Type:         errorType,
		FromCompiler: false,
	}
}

var (
	ErrRuntime          = errors.New(RunTimeErrMsg)
	ErrStringLimit      = errors.New("string Limit")
	ErrCompilation      = errors.New(CompilerErrorMsg)
	ErrAssertionFailure = errors.New(AssertionErr)
)
