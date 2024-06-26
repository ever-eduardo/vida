package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/lexer"
	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

var compilerError verror.VidaError

type Compiler struct {
	locals          locals
	current         tokenInfo
	next            tokenInfo
	compilationInfo *compilationInfo
	lexer           *lexer.Lexer
	parent          *Compiler
	module          *Module
	function        *Function
	identifiersMap  map[string]uint16
	kIndex          uint16
	rA              byte
	rB              byte
	rC              byte
	ok              bool
}

func NewCompiler(src []byte, moduleName string) *Compiler {
	c := Compiler{
		lexer:           lexer.New(src, moduleName),
		identifiersMap:  make(map[string]uint16),
		module:          NewModule(),
		parent:          nil,
		locals:          locals{},
		ok:              true,
		compilationInfo: &compilationInfo{},
	}
	c.advance()
	c.advance()
	c.makeHeader()
	return &c
}

func newChildCompiler(c *Compiler) Compiler {
	child := Compiler{
		lexer:           c.lexer,
		identifiersMap:  c.identifiersMap,
		function:        &Function{},
		module:          c.module,
		parent:          c,
		locals:          locals{},
		ok:              true,
		compilationInfo: c.compilationInfo,
	}
	return child
}

func (c *Compiler) Compile() (*Module, error) {
	for c.ok {
		switch c.current.token {
		case token.IDENTIFIER:
			c.identifierPath()
		case token.LOCAL:
			c.localDecl()
		case token.EOF:
			c.makeStopRun()
			return c.module, nil
		default:
			return nil, verror.New(c.lexer.ModuleName, "Expected statement", verror.SyntaxError, c.current.line)
		}
	}
	return nil, compilerError
}

func (c *Compiler) identifierPath() {
	c.compilationInfo.Identifier = c.current.lit
	c.compilationInfo.IsGlobalAssignment = true
	c.advance()
	c.expect(token.ASSIGN)
	c.advance()
	c.expression()
	c.advance()
	c.makeIdentifierPath()
}

func (c *Compiler) localDecl() {
	c.advance()
	c.expect(token.IDENTIFIER)
	id := c.current.lit
	c.advance()
	c.expression()
	c.advance()
	c.makeLocal(id)
}

func (c *Compiler) expression() {
	switch c.current.token {
	case token.TRUE, token.FALSE, token.NIL:
		c.makeAtomic(c.current.token)
	case token.IDENTIFIER:
		c.makeLoadGlobal()
	default:
	}
}

func (c *Compiler) expect(tok token.Token) {
	if c.current.token != tok && c.ok {
		c.ok = false
		message := fmt.Sprintf("Expected token %v and got token %v", tok, c.current.token)
		compilerError = verror.New(c.lexer.ModuleName, message, verror.SyntaxError, c.current.line)
	}
}

func (c *Compiler) advance() token.Token {
	c.current.line, c.current.token, c.current.lit = c.next.line, c.next.token, c.next.lit
	c.next.line, c.next.token, c.next.lit = c.lexer.Next()
	return c.current.token
}

type tokenInfo struct {
	lit   string
	token token.Token
	line  uint
}

type localInfo struct {
	Name     string
	Register int
}

type locals struct {
	Locals []localInfo
}

func (local locals) add(name string, register int) {
	local.Locals = append(local.Locals, localInfo{Name: name, Register: register})
}

func (local locals) lastIndexOf(name string) (int, bool) {
	for i := len(local.Locals) - 1; i >= 0; i-- {
		if local.Locals[i].Name == name {
			return local.Locals[i].Register, true
		}
	}
	return 0, false
}

type compilationInfo struct {
	Identifier         string
	IsGlobalAssignment bool
	IsAtomicAssignment bool
}
