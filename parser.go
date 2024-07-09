package vida

import (
	"fmt"
	"strconv"

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
		case token.IDENTIFIER:
			p.ast.Statement = append(p.ast.Statement, p.identPath())
		case token.LOC:
			p.ast.Statement = append(p.ast.Statement, p.localStmt())
		case token.LCURLY:
			p.ast.Statement = append(p.ast.Statement, p.block())
		case token.EOF:
			return p.ast, nil
		default:
			p.err = verror.New(p.lexer.ModuleName, "Expected statement", verror.SyntaxErrMsg, p.current.Line)
			return nil, p.err
		}
	}
	return nil, p.err
}

func (p *Parser) identPath() ast.Node {
	i := p.current.Lit
	p.advance()
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression(token.LowestPrec)
	p.advance()
	return &ast.Set{LHS: &ast.Identifier{Value: i}, Expr: e}
}

func (p *Parser) localStmt() ast.Node {
	p.advance()
	p.expect(token.IDENTIFIER)
	i := p.current.Lit
	p.advance()
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression(token.LowestPrec)
	p.advance()
	return &ast.Loc{Identifier: i, Expr: e}
}

func (p *Parser) block() ast.Node {
	block := &ast.Block{}
	p.advance()
	for p.current.Token != token.RCURLY {
		switch p.current.Token {
		case token.IDENTIFIER:
			block.Statement = append(block.Statement, p.identPath())
		case token.LOC:
			block.Statement = append(block.Statement, p.localStmt())
		case token.LCURLY:
			block.Statement = append(block.Statement, p.block())
		default:
			p.err = verror.New(p.lexer.ModuleName, "Expected statement", verror.SyntaxErrMsg, p.current.Line)
			return block
		}
	}
	p.advance()
	return block
}

func (p *Parser) expression(precedence int) ast.Node {
	e := p.prefix()
	for p.next.Token.IsBinaryOperator() && p.next.Token.Precedence() > precedence {
		p.advance()
		op := p.current.Token
		p.advance()
		r := p.expression(op.Precedence())
		e = &ast.BinaryExpr{Op: op, Lhs: e, Rhs: r}
	}
	return e
}

func (p *Parser) prefix() ast.Node {
	switch p.current.Token {
	case token.NOT, token.SUB, token.ADD:
		t := p.current.Token
		p.advance()
		e := p.prefix()
		return &ast.PrefixExpr{Op: t, Expr: e}
	}
	return p.primary()
}

func (p *Parser) primary() ast.Node {
	e := p.operand()
Loop:
	for p.next.Token == token.LBRACKET {
		p.advance()
		switch p.current.Token {
		case token.LBRACKET:
			e = p.indexOrSlice(e)
		default:
			break Loop
		}
	}
	return e
}

func (p *Parser) operand() ast.Node {
	switch p.current.Token {
	case token.INTEGER:
		if i, err := strconv.ParseInt(p.current.Lit, 0, 64); err == nil {
			return &ast.Integer{Value: i}
		}
		p.err = verror.New(p.lexer.ModuleName, "Integer value could not be processed", verror.SyntaxErrMsg, p.current.Line)
		p.ok = false
		return &ast.Nil{}
	case token.FLOAT:
		if f, err := strconv.ParseFloat(p.current.Lit, 64); err == nil {
			return &ast.Float{Value: f}
		}
		p.err = verror.New(p.lexer.ModuleName, "Float value could not be processed", verror.SyntaxErrMsg, p.current.Line)
		p.ok = false
		return &ast.Nil{}
	case token.STRING:
		return &ast.String{Value: p.current.Lit}
	case token.TRUE:
		return &ast.Boolean{Value: true}
	case token.FALSE:
		return &ast.Boolean{Value: false}
	case token.NIL:
		return &ast.Nil{}
	case token.IDENTIFIER:
		return &ast.Reference{Value: p.current.Lit}
	case token.LBRACKET:
		xs := &ast.List{}
		p.advance()
		for p.current.Token != token.RBRACKET {
			e := p.expression(token.LowestPrec)
			p.advance()
			xs.ExprList = append(xs.ExprList, e)
			for p.current.Token == token.COMMA {
				p.advance()
				if p.current.Token == token.RBRACKET {
					p.expect(token.RBRACKET)
					return xs
				}
				e := p.expression(token.LowestPrec)
				p.advance()
				xs.ExprList = append(xs.ExprList, e)
			}
		}
		p.expect(token.RBRACKET)
		return xs
	case token.LPAREN:
		p.advance()
		if p.current.Token == token.RPAREN {
			p.err = verror.New(p.lexer.ModuleName, "Expected expression", verror.SyntaxErrMsg, p.current.Line)
			p.ok = false
			return &ast.Nil{}
		}
		e := p.expression(token.LowestPrec)
		p.advance()
		p.expect(token.RPAREN)
		return e
	default:
		p.err = verror.New(p.lexer.ModuleName, "Expected expression", verror.SyntaxErrMsg, p.current.Line)
		p.ok = false
		return &ast.Nil{}
	}
}

func (p *Parser) indexOrSlice(e ast.Node) ast.Node {
	p.advance()
	var index [2]ast.Node
	if p.current.Token != token.COLON {
		index[0] = p.expression(token.LowestPrec)
		p.advance()
	}
	var numColons int
	if p.current.Token == token.COLON {
		numColons++
		p.advance()
		if p.current.Token != token.RBRACKET && p.current.Token != token.EOF {
			index[1] = p.expression(token.LowestPrec)
			p.advance()
		}
	}
	p.expect(token.RBRACKET)
	if numColons > 0 {
		return &ast.SliceGet{
			List:  e,
			First: index[0],
			Last:  index[1],
		}
	}
	return &ast.IndexGet{
		Indexable: e,
		Index:     index[0],
	}
}

func (p *Parser) expect(tok token.Token) {
	if p.current.Token != tok && p.ok {
		p.ok = false
		message := fmt.Sprintf("Expected token %v and got token %v", tok, p.current.Token)
		p.err = verror.New(p.lexer.ModuleName, message, verror.SyntaxErrMsg, p.current.Line)
	}
}

func (p *Parser) advance() token.Token {
	p.current.Line, p.current.Token, p.current.Lit = p.next.Line, p.next.Token, p.next.Lit
	p.next.Line, p.next.Token, p.next.Lit = p.lexer.Next()
	return p.current.Token
}
