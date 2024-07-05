package verror

import "fmt"

const FileError = "File"
const LexicalError = "Lexical"
const SyntaxError = "Syntax"
const RunTimeError = "Runtime"

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
