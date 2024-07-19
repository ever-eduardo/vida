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
			p.ast.Statement = append(p.ast.Statement, p.identPath(&p.ast.Statement))
		case token.LOC:
			p.ast.Statement = append(p.ast.Statement, p.localStmt())
		case token.IF:
			p.ast.Statement = append(p.ast.Statement, p.ifStmt(false))
		case token.FOR:
			p.ast.Statement = append(p.ast.Statement, p.forLoop())
		case token.WHILE:
			p.ast.Statement = append(p.ast.Statement, p.loop())
		case token.LCURLY:
			p.ast.Statement = append(p.ast.Statement, p.block(false))
			p.advance()
		case token.EOF:
			return p.ast, nil
		default:
			p.err = verror.New(p.lexer.ModuleName, "Expected statement", verror.SyntaxErrMsg, p.current.Line)
			return nil, p.err
		}
	}
	return nil, p.err
}

func (p *Parser) identPath(statements *[]ast.Node) ast.Node {
	if p.next.Token == token.DOT || p.next.Token == token.LBRACKET {
		return p.mutateDataStructureOrCall(statements)
	}
	if p.next.Token == token.LPAREN {
		return p.callStmt(statements)
	}
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

func (p *Parser) block(isInsideLoop bool) ast.Node {
	block := &ast.Block{}
	p.advance()
	for p.current.Token != token.RCURLY {
		switch p.current.Token {
		case token.IDENTIFIER:
			block.Statement = append(block.Statement, p.identPath(&block.Statement))
		case token.LOC:
			block.Statement = append(block.Statement, p.localStmt())
		case token.IF:
			block.Statement = append(block.Statement, p.ifStmt(isInsideLoop))
		case token.FOR:
			block.Statement = append(block.Statement, p.forLoop())
		case token.WHILE:
			block.Statement = append(block.Statement, p.loop())
		case token.RET:
			block.Statement = append(block.Statement, p.ret())
		case token.BREAK:
			if isInsideLoop {
				block.Statement = append(block.Statement, p.breakStmt())
			} else {
				p.err = verror.New(p.lexer.ModuleName, "Break outside loop", verror.SyntaxErrMsg, p.current.Line)
				p.ok = false
				return block
			}
		case token.CONTINUE:
			if isInsideLoop {
				block.Statement = append(block.Statement, p.continueStmt())
			} else {
				p.err = verror.New(p.lexer.ModuleName, "Continue outside loop", verror.SyntaxErrMsg, p.current.Line)
				p.ok = false
				return block
			}
		case token.LCURLY:
			block.Statement = append(block.Statement, p.block(isInsideLoop))
			p.advance()
		default:
			p.err = verror.New(p.lexer.ModuleName, "Expected statement", verror.SyntaxErrMsg, p.current.Line)
			p.ok = false
			return block
		}
	}
	return block
}

func (p *Parser) mutateDataStructureOrCall(statements *[]ast.Node) ast.Node {
	*statements = append(*statements, &ast.ReferenceStmt{Value: p.current.Lit})
	var i ast.Node
Loop:
	for p.next.Token == token.LBRACKET || p.next.Token == token.DOT {
		p.advance()
		switch p.current.Token {
		case token.LBRACKET:
			p.advance()
			i = p.expression(token.LowestPrec)
			p.advance()
			p.expect(token.RBRACKET)
			if p.next.Token == token.ASSIGN {
				break Loop
			}
			*statements = append(*statements, &ast.IGetStmt{Index: i})
		case token.DOT:
			p.advance()
			p.expect(token.IDENTIFIER)
			i = &ast.Property{Value: p.current.Lit}
			if p.next.Token == token.ASSIGN {
				break Loop
			}
			*statements = append(*statements, &ast.SelectStmt{Selector: i})
		default:
			break Loop
		}
	}
	p.advance()
	if p.current.Token == token.LPAREN {
		p.advance()
		var args []ast.Node
		for p.current.Token != token.RPAREN && p.current.Token != token.EOF {
			args = append(args, p.expression(token.LowestPrec))
			p.advance()
			if p.current.Token == token.COMMA {
				for p.current.Token == token.COMMA {
					p.advance()
					args = append(args, p.expression(token.LowestPrec))
					p.advance()
				}
				goto afterParen
			}
		}
	afterParen:
		p.expect(token.RPAREN)
		p.advance()
		return &ast.CallStmt{Args: args}
	}
	p.expect(token.ASSIGN)
	p.advance()
	e := p.expression(token.LowestPrec)
	p.advance()
	return &ast.ISet{Index: i, Expr: e}
}

func (p *Parser) callStmt(statements *[]ast.Node) ast.Node {
	*statements = append(*statements, &ast.ReferenceStmt{Value: p.current.Lit})
	var args []ast.Node
	p.advance()
	p.expect(token.LPAREN)
	p.advance()
	for p.current.Token != token.RPAREN && p.current.Token != token.EOF {
		args = append(args, p.expression(token.LowestPrec))
		p.advance()
		if p.current.Token == token.COMMA {
			for p.current.Token == token.COMMA {
				p.advance()
				args = append(args, p.expression(token.LowestPrec))
				p.advance()
			}
			goto afterParen
		}
	}
afterParen:
	p.expect(token.RPAREN)
	p.advance()
	return &ast.CallStmt{Args: args}
}

func (p *Parser) forLoop() ast.Node {
	p.advance()
	p.expect(token.IDENTIFIER)
	id := p.current.Lit
	p.advance()
	if p.current.Token == token.COMMA {
		return p.iterforLoop(id)
	}
	p.expect(token.ASSIGN)
	p.advance()
	init := p.expression(token.LowestPrec)
	p.advance()
	p.expect(token.COMMA)
	p.advance()
	end := p.expression(token.LowestPrec)
	p.advance()
	var step ast.Node = &ast.Integer{Value: 1}
	if p.current.Token == token.COMMA {
		p.expect(token.COMMA)
		p.advance()
		step = p.expression(token.LowestPrec)
		p.advance()
	}
	p.expect(token.LCURLY)
	block := p.block(true)
	p.advance()
	return &ast.For{Init: init, End: end, Id: id, Step: step, Block: block}
}

func (p *Parser) iterforLoop(key string) ast.Node {
	p.advance()
	p.expect(token.IDENTIFIER)
	v := p.current.Lit
	p.advance()
	p.expect(token.IN)
	p.advance()
	e := p.expression(token.LowestPrec)
	p.advance()
	p.expect(token.LCURLY)
	b := p.block(true)
	p.advance()
	return &ast.IFor{Key: key, Value: v, Expr: e, Block: b}
}

func (p *Parser) ifStmt(isInsideLoop bool) ast.Node {
	p.advance()
	c := p.expression(token.LowestPrec)
	p.advance()
	p.expect(token.LCURLY)
	b := p.block(isInsideLoop)
	p.advance()
	branch := &ast.Branch{If: &ast.If{Condition: c, Block: b}}
	for p.current.Token == token.ELSE && p.next.Token == token.IF {
		p.advance()
		p.advance()
		c := p.expression(token.LowestPrec)
		p.advance()
		p.expect(token.LCURLY)
		b := p.block(isInsideLoop)
		p.advance()
		branch.Elifs = append(branch.Elifs, &ast.If{Condition: c, Block: b})
	}
	if p.current.Token == token.ELSE {
		p.advance()
		b := p.block(isInsideLoop)
		p.advance()
		branch.Else = &ast.Else{Block: b}
	}
	return branch
}

func (p *Parser) loop() ast.Node {
	p.advance()
	c := p.expression(token.LowestPrec)
	p.advance()
	p.expect(token.LCURLY)
	b := p.block(true)
	p.advance()
	return &ast.While{Condition: c, Block: b}
}

func (p *Parser) breakStmt() ast.Node {
	p.advance()
	return &ast.Break{}
}

func (p *Parser) continueStmt() ast.Node {
	p.advance()
	return &ast.Continue{}
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
	for p.next.Token == token.LBRACKET || p.next.Token == token.DOT || p.next.Token == token.LPAREN {
		p.advance()
		switch p.current.Token {
		case token.LBRACKET:
			e = p.indexOrSlice(e)
		case token.DOT:
			p.advance()
			switch p.current.Token {
			case token.IDENTIFIER:
				e = p.selector(e)
			default:
				p.err = verror.New(p.lexer.ModuleName, "Expected identifier", verror.SyntaxErrMsg, p.current.Line)
				p.ok = false
				return &ast.Nil{}
			}
		case token.LPAREN:
			e = p.callExpr(e)
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
			if p.current.Token == token.COMMA {
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
			goto endList
		}
	endList:
		p.expect(token.RBRACKET)
		return xs
	case token.LCURLY:
		doc := &ast.Document{}
		p.advance()
		for p.current.Token != token.RCURLY {
			p.expect(token.IDENTIFIER)
			k := &ast.Property{Value: p.current.Lit}
			p.advance()
			p.expect(token.COLON)
			p.advance()
			v := p.expression(token.LowestPrec)
			p.advance()
			doc.Pairs = append(doc.Pairs, &ast.Pair{Key: k, Value: v})
			if p.current.Token == token.COMMA {
				for p.current.Token == token.COMMA {
					p.advance()
					if p.current.Token == token.RCURLY {
						p.expect(token.RCURLY)
						return doc
					}
					p.expect(token.IDENTIFIER)
					k := &ast.Property{Value: p.current.Lit}
					p.advance()
					p.expect(token.COLON)
					p.advance()
					v := p.expression(token.LowestPrec)
					p.advance()
					doc.Pairs = append(doc.Pairs, &ast.Pair{Key: k, Value: v})
				}
			}
			goto endDoc
		}
	endDoc:
		p.expect(token.RCURLY)
		return doc
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
	case token.FUN:
		f := &ast.Fun{}
		p.advance()
		for p.current.Token == token.IDENTIFIER {
			f.Args = append(f.Args, p.current.Lit)
			p.advance()
			if p.current.Token == token.COMMA {
				for p.current.Token == token.COMMA {
					p.advance()
					p.expect(token.IDENTIFIER)
					f.Args = append(f.Args, p.current.Lit)
					p.advance()
				}
			}
			goto endParams
		}
	endParams:
		if p.current.Token == token.ARROW {
			p.advance()
			e := p.expression(token.LowestPrec)
			b := &ast.Block{}
			b.Statement = append(b.Statement, &ast.Ret{Expr: e})
			f.Body = b
			return f
		}
		p.expect(token.LCURLY)
		block := p.block(false)
		block.(*ast.Block).Statement = append(block.(*ast.Block).Statement, &ast.Ret{Expr: &ast.Nil{}})
		f.Body = block
		return f
	default:
		p.err = verror.New(p.lexer.ModuleName, "Expected expression", verror.SyntaxErrMsg, p.current.Line)
		p.ok = false
		return &ast.Nil{}
	}
}

func (p *Parser) ret() ast.Node {
	p.advance()
	e := p.expression(token.LowestPrec)
	p.advance()
	return &ast.Ret{Expr: e}
}

func (p *Parser) callExpr(e ast.Node) ast.Node {
	p.advance()
	var args []ast.Node
	for p.current.Token != token.RPAREN && p.current.Token != token.EOF {
		args = append(args, p.expression(token.LowestPrec))
		p.advance()
		if p.current.Token == token.COMMA {
			for p.current.Token == token.COMMA {
				p.advance()
				args = append(args, p.expression(token.LowestPrec))
				p.advance()
			}
			goto afterParen
		}
	}
afterParen:
	p.expect(token.RPAREN)
	return &ast.CallExpr{Fun: e, Args: args}
}

func (p *Parser) indexOrSlice(e ast.Node) ast.Node {
	p.advance()
	var index [2]ast.Node
	mode := 2
	if p.current.Token != token.COLON {
		mode |= 4
		index[0] = p.expression(token.LowestPrec)
		p.advance()
	}
	var numColons int
	if p.current.Token == token.COLON {
		numColons++
		p.advance()
		if p.current.Token != token.RBRACKET && p.current.Token != token.EOF {
			mode |= 1
			index[1] = p.expression(token.LowestPrec)
			p.advance()
		}
	}
	p.expect(token.RBRACKET)
	if numColons > 0 {
		return &ast.Slice{
			Value: e,
			First: index[0],
			Last:  index[1],
			Mode:  mode,
		}
	}
	return &ast.IGet{
		Indexable: e,
		Index:     index[0],
	}
}

func (p *Parser) selector(e ast.Node) ast.Node {
	return &ast.Select{Selectable: e, Selector: &ast.Property{Value: p.current.Lit}}
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
