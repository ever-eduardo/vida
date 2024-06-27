package vida

import (
	"github.com/ever-eduardo/vida/lexer"
	"github.com/ever-eduardo/vida/token"
)

type Compiler struct {
	lexer    lexer.Lexer
	module   Module
	current  token.TokenInfo
	next     token.TokenInfo
	level    int
	parent   *Compiler
	function *Function
}
