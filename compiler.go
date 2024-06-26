package vida

import "github.com/ever-eduardo/vida/ast"

type Compiler struct {
	ast      *ast.Ast
	module   *Module
	function *Function
	parent   *Compiler
	kb       *KonstBuilder
	sb       *symbolBuilder
	scope    int
	level    int
	rAlloc   byte
}

func NewCompiler(ast *ast.Ast, moduleName string) *Compiler {
	return &Compiler{
		ast:    ast,
		module: newModule(moduleName),
		kb:     newKonstBuilder(),
		sb:     newSymbolBuilder(),
	}
}

func newChildCompiler(p *Compiler) *Compiler {
	return &Compiler{
		ast:      p.ast,
		module:   p.module,
		function: newFunction(),
		kb:       p.kb,
		sb:       p.sb,
		parent:   p,
		level:    p.level + 1,
	}
}

func (c *Compiler) CompileModule() *Module {
	c.appendHeader()
	for i := range len(c.ast.Statement) {
		c.compileStmt(c.ast.Statement[i])
	}
	c.module.Konstants = c.kb.Konstants
	c.appendEnd()
	return c.module
}

func (c *Compiler) compileStmt(node ast.Node) {
	switch n := node.(type) {
	case *ast.Set:
		from, scope := c.compileExpr(n.Expr)
		switch lhs := n.LHS.(type) {
		case *ast.Identifier:
			if to, isLocal := c.sb.isLocal(lhs.Value, c.level, c.scope); isLocal {
				if scope == rLocal {
					c.emitMove(byte(from), to)
				} else {
					c.emitLoc(from, to, scope)
				}
			} else if isGlobal := c.sb.isGlobal(lhs.Value); isGlobal {
				to := c.kb.StringIndex(lhs.Value)
				c.emitSetG(from, to, scope)
			} else {
				to := c.kb.StringIndex(lhs.Value)
				c.sb.addGlobal(lhs.Value)
				c.emitSetG(from, to, scope)
			}
		}
	case *ast.Loc:
		to := c.rAlloc
		c.rAlloc++
		from, scope := c.compileExpr(n.Expr)
		c.sb.addLocal(n.Identifier, c.level, c.scope, to)
		c.emitLoc(from, to, scope)
	}
}

func (c *Compiler) compileExpr(node ast.Node) (int, byte) {
	switch n := node.(type) {
	case *ast.PrefixExpr:
		idx, scope := c.compileExpr(n.Expr)
		c.emitNot(idx, c.rAlloc, scope)
		return int(c.rAlloc), rLocal
	case *ast.Boolean:
		idx := c.kb.BooleanIndex(n.Value)
		return idx, rKonst
	case *ast.Nil:
		idx := c.kb.NilIndex()
		return idx, rKonst
	case *ast.Reference:
		idx, scope := c.refScope(n.Value)
		return idx, scope
	default:
		return 0, rKonst
	}
}
