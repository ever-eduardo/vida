package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/ast"
	"github.com/ever-eduardo/vida/lexer"
	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Parser struct {
	err     verror.VidaError
	current token.TokenInfo
	next    token.TokenInfo
	lexer   *lexer.Lexer
	ast     *ast.Ast
	ok      bool
}

func NewParser(src []byte, moduleName string) *Parser {
	p := Parser{
		lexer: lexer.New(src, moduleName),
		ok:    true,
		ast:   &ast.Ast{},
	}
	p.advance()
	p.advance()
	return &p
}

func (p *Parser) Parse() (*ast.Ast, error) {
	for p.ok {
		switch p.current.Token {
		case token.VAL:
			p.ast.Statement = append(p.ast.Statement, p.globalDecl())
		case token.LOC:
			p.ast.Statement = append(p.ast.Statement, p.localDecl())
		case token.IDENTIFIER:
			p.ast.Statement = append(p.ast.Statement, p.mutState())
		case token.EOF:
			return p.ast, nil
		default:
			p.err = verror.New(p.lexer.ModuleName, "Expected statement", verror.SyntaxError, p.current.Line)
			return nil, p.err
		}
	}
	return p.ast, nil
}

func (p *Parser) globalDecl() ast.Statement {
	p.advance()
	p.expect(token.IDENTIFIER)
	i := p.current.Lit
	p.advance()
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression()
	p.advance()
	return ast.Val{Identifier: i, Expr: e}
}

func (p *Parser) localDecl() ast.Statement {
	p.advance()
	p.expect(token.IDENTIFIER)
	i := p.current.Lit
	p.advance()
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression()
	p.advance()
	return ast.Loc{Identifier: i, Expr: e}
}

func (p *Parser) mutState() ast.Statement {
	p.expect(token.IDENTIFIER)
	lhs := p.current.Lit
	p.advance()
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression()
	p.advance()
	return ast.Mut{Identifier: lhs, Expr: e}
}

func (p *Parser) expression() ast.Expr {
	switch p.current.Token {
	case token.TRUE:
		return ast.Boolean{Value: true}
	case token.FALSE:
		return ast.Boolean{Value: false}
	case token.NIL:
		return ast.Nil{}
	default:
		return ast.Reference{Identifier: p.current.Lit}
	}
}

func (p *Parser) expect(tok token.Token) {
	if p.current.Token != tok && p.ok {
		p.ok = false
		message := fmt.Sprintf("Expected token %v and got token %v", tok, p.current.Token)
		p.err = verror.New(p.lexer.ModuleName, message, verror.SyntaxError, p.current.Line)
	}
}

func (p *Parser) advance() token.Token {
	p.current.Line, p.current.Token, p.current.Lit = p.next.Line, p.next.Token, p.next.Lit
	p.next.Line, p.next.Token, p.next.Lit = p.lexer.Next()
	return p.current.Token
}
