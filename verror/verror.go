package verror

import (
	"errors"
	"fmt"
)

const FileErrMsg = "File"
const LexicalErrMsg = "Lexical"
const SyntaxErrMsg = "Syntax"
const CompilerErrorMsg = "Compiler error"
const RunTimeErrMsg = "Runtime error"
const AssertionErr = "Assertion failure"
const Exception = "An exception was raised"

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
	ErrStringLimit      = errors.New("Reached string limit")
	ErrCompilation      = errors.New(CompilerErrorMsg)
	ErrAssertionFailure = errors.New(AssertionErr)
)
