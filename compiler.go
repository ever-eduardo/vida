package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/lexer"
	"github.com/ever-eduardo/vida/symbol"
	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Compiler struct {
	lexer           lexer.Lexer
	module          Module
	current         token.TokenInfo
	next            token.TokenInfo
	symbolTable     symbol.SymbolTable
	compilationInfo compilationInfo
	level           int
	parent          *Compiler
	function        *Function
	kIndex          uint16
	rA              byte
	rB              byte
	rC              byte
	ok              bool
}

func NewCompiler(src []byte, moduleName string) *Compiler {
	c := Compiler{
		lexer:           lexer.New(src, moduleName),
		symbolTable:     symbol.New(),
		module:          NewModule(moduleName),
		parent:          nil,
		ok:              true,
		compilationInfo: compilationInfo{},
	}
	c.advance()
	c.advance()
	c.makeHeader()
	return &c
}

func newChildCompiler(c *Compiler) Compiler {
	child := Compiler{
		lexer:           c.lexer,
		symbolTable:     c.symbolTable,
		function:        &Function{},
		module:          c.module,
		parent:          c,
		ok:              true,
		compilationInfo: c.compilationInfo,
	}
	return child
}

func (c *Compiler) Compile() (Module, error) {
	for c.ok {
		switch c.current.Token {
		case token.LET:
			c.defineGlobal()
		case token.VAR:
			c.defineLocal()
		case token.IDENTIFIER:
			c.changeStateOrCall()
		case token.EOF:
			c.makeStopRun()
			return c.module, nil
		default:
			return c.module, verror.New(c.lexer.ModuleName, "Expected valid statement", verror.SyntaxError, c.current.Line)
		}
	}
	return c.module, compilerError
}

func (c *Compiler) defineGlobal() {
	c.compilationInfo.Let = true
	c.advance()
	c.expect(token.IDENTIFIER)
	c.compilationInfo.Identifier = c.current.Lit
	c.advance()
	c.expect(token.ASSIGN)
	c.advance()
	c.expression()
	c.advance()
	c.makeNewGlobal()
}

func (c *Compiler) defineLocal() {
	c.compilationInfo.IsLocalAssignment = true
	c.advance()
	c.expect(token.IDENTIFIER)
	c.compilationInfo.Identifier = c.current.Lit
	c.advance()
	c.expect(token.ASSIGN)
	c.advance()
	c.expression()
	c.advance()
	c.makeLocal()
}

func (c *Compiler) changeStateOrCall() {
	c.compilationInfo.Identifier = c.current.Lit
}

func (c *Compiler) expression() {
	switch c.current.Token {
	case token.TRUE, token.FALSE, token.NIL:
		c.makeAtomic(c.current.Token)
	case token.IDENTIFIER:
		c.makeLoadRef()
	default:
		c.ok = false
		message := fmt.Sprintf("Expected expression and got token %v", c.current.Lit)
		compilerError = verror.New(c.lexer.ModuleName, message, verror.SyntaxError, c.current.Line)
	}
}

func (c *Compiler) expect(tok token.Token) {
	if c.current.Token != tok && c.ok {
		c.ok = false
		message := fmt.Sprintf("Expected token %v and got token %v", tok, c.current.Token)
		compilerError = verror.New(c.lexer.ModuleName, message, verror.SyntaxError, c.current.Line)
	}
}

func (c *Compiler) advance() token.Token {
	c.current.Line, c.current.Token, c.current.Lit = c.next.Line, c.next.Token, c.next.Lit
	c.next.Line, c.next.Token, c.next.Lit = c.lexer.Next()
	return c.current.Token
}

type compilationInfo struct {
	Identifier         string
	Let                bool
	Var                bool
	IsGlobalAssignment bool
	IsLocalAssignment  bool
	IsAtom             bool
}
