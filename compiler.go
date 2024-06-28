package vida

import "github.com/ever-eduardo/vida/ast"

type Compiler struct {
	ast      *ast.Ast
	module   *Module
	function *Function
	parent   *Compiler
	kb       *KonstBuilder
	scope    int
	fromK    bool
}

func NewCompiler(ast *ast.Ast, moduleName string) *Compiler {
	return &Compiler{
		ast:    ast,
		module: newModule(moduleName),
		kb:     newKonstBuilder(),
	}
}

func newChildCompiler(p *Compiler) *Compiler {
	return &Compiler{
		ast:      p.ast,
		module:   p.module,
		function: newFunction(),
		kb:       p.kb,
		parent:   p,
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
			dest := c.kb.StringIndex(lhs.Value)
			if src, flag, isKS := c.compileExpr(n.Expr); isKS {
				c.emitSetFromK(dest, src, flag)
			}
		}
	}
}

func (c *Compiler) compileExpr(node ast.Node) (int, byte, bool) {
	switch n := node.(type) {
	case *ast.Boolean:
		idx := c.kb.BooleanIndex(n.Value)
		return idx, 0, true
	case *ast.Nil:
		idx := c.kb.NilIndex()
		return idx, 0, true
	case *ast.Reference:
		idx := c.kb.StringIndex(n.Value)
		return idx, 1, true
	default:
		return 0, 0, false
	}
}
