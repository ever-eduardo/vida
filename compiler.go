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
			if to, isLocal := c.sb.isLocal(lhs.Value); isLocal {
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
	case *ast.Block:
		c.scope++
		for i := range len(n.Statement) {
			c.compileStmt(n.Statement[i])
		}
		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.scope--
	}
}

func (c *Compiler) compileExpr(node ast.Node) (int, byte) {
	switch n := node.(type) {
	case *ast.Integer:
		idx := c.kb.IntegerIndex(n.Value)
		return idx, rKonst
	case *ast.Float:
		idx := c.kb.FloatIndex(n.Value)
		return idx, rKonst
	case *ast.String:
		idx := c.kb.StringIndex(n.Value)
		return idx, rKonst
	case *ast.BinaryExpr:
		opReg := c.rAlloc
		c.rAlloc++
		lidx, lscope := c.compileExpr(n.Lhs)
		c.rAlloc++
		ridx, rscope := c.compileExpr(n.Rhs)
		c.rAlloc = opReg
		c.emitBinary(lidx, ridx, lscope, rscope, opReg, byte(n.Op))
		return int(opReg), rLocal
	case *ast.PrefixExpr:
		idx, scope := c.compileExpr(n.Expr)
		c.emitPrefix(idx, c.rAlloc, scope, byte(n.Op))
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
	case *ast.List:
		if len(n.ExprList) == 0 {
			c.emitList(0, c.rAlloc, c.rAlloc)
			return int(c.rAlloc), rLocal
		}
		var count int
		for _, v := range n.ExprList {
			idx, scope := c.compileExpr(v)
			if scope != rLocal {
				c.emitLoc(idx, c.rAlloc, scope)
			} else if idx != int(c.rAlloc) {
				c.emitMove(byte(idx), c.rAlloc)
			}
			c.rAlloc++
			count++
		}
		c.rAlloc -= byte(count)
		c.emitList(byte(count), c.rAlloc, c.rAlloc)
		return int(c.rAlloc), rLocal
	case *ast.IndexGet:
		resultReg := c.rAlloc
		c.rAlloc++
		i, s := c.compileExpr(n.Indexable)
		c.rAlloc++
		j, t := c.compileExpr(n.Index)
		c.rAlloc = resultReg
		c.emitIndexGet(i, j, s, t, resultReg)
		return int(resultReg), rLocal
	case *ast.Slice:
		resultReg := c.rAlloc
		var scopeV, scopeL, scopeR byte
		var fromV, fromL, fromR int
		c.rAlloc++
		fromV, scopeV = c.compileExpr(n.Value)
		switch n.Mode {
		case vcv:
			break
		case vce:
			c.rAlloc++
			fromR, scopeR = c.compileExpr(n.Last)
		case ecv:
			c.rAlloc++
			fromL, scopeL = c.compileExpr(n.First)
		case ece:
			c.rAlloc++
			fromL, scopeL = c.compileExpr(n.First)
			c.rAlloc++
			fromR, scopeR = c.compileExpr(n.Last)
		}
		c.rAlloc = resultReg
		c.emitSlice(byte(n.Mode), fromV, fromL, fromR, scopeV, scopeL, scopeR, resultReg)
		return int(resultReg), rLocal
	default:
		return 0, rKonst
	}
}
