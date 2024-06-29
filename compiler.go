package vida

import "github.com/ever-eduardo/vida/ast"

type Compiler struct {
	ast      *ast.Ast
	module   *Module
	function *Function
	parent   *Compiler
	kb       *KonstBuilder
	lb       *LocalBuilder
	scope    int
	level    int
	fromK    bool
	lrAlloc  byte
}

func NewCompiler(ast *ast.Ast, moduleName string) *Compiler {
	return &Compiler{
		ast:    ast,
		module: newModule(moduleName),
		kb:     newKonstBuilder(),
		lb:     NewLocalBuilder(),
	}
}

func newChildCompiler(p *Compiler) *Compiler {
	return &Compiler{
		ast:      p.ast,
		module:   p.module,
		function: newFunction(),
		kb:       p.kb,
		lb:       p.lb,
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
		switch lhs := n.LHS.(type) {
		case *ast.Identifier:
			if to, isLocal := c.lb.IsLocal(lhs.Value, c.level, c.scope); isLocal {
				if from, flag, isSK := c.compileExpr(n.Expr); isSK {
					c.emitLocSK(from, to, flag)
				} else {
					c.emitMove(byte(from), to)
				}
			} else if to, isGlobal := c.kb.IsGlobal(lhs.Value); isGlobal {
				from, flag, _ := c.compileExpr(n.Expr)
				c.emitSetSK(from, to, flag)
			} else {
				to := c.kb.StringIndex(lhs.Value)
				from, flag, _ := c.compileExpr(n.Expr)
				c.emitSetSK(from, to, flag)
			}
		}
	case *ast.Loc:
		to := c.lrAlloc
		c.lb.AddLocal(n.Identifier, c.level, c.scope, c.lrAlloc)
		c.lrAlloc++
		if from, flag, isSK := c.compileExpr(n.Expr); isSK {
			c.emitLocSK(from, to, flag)
		} else {
			c.emitMove(byte(from), to)
		}
	}
}

func (c *Compiler) compileExpr(node ast.Node) (int, byte, bool) {
	switch n := node.(type) {
	case *ast.Boolean:
		idx := c.kb.BooleanIndex(n.Value)
		return idx, refKns, true
	case *ast.Nil:
		idx := c.kb.NilIndex()
		return idx, refKns, true
	case *ast.Reference:
		idx, s := c.referenceScope(n.Value)
		switch s {
		case refReg:
			return idx, s, false
		case refStr:
			return idx, s, true
		default:
			return idx, s, false
		}
	default:
		return 0, 0, false
	}
}
