package vida

import (
	"encoding/binary"

	"github.com/ever-eduardo/vida/ast"
	"github.com/ever-eduardo/vida/token"
)

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
	case *ast.For:
		c.scope++
		var reg [4]byte

		initIdx, initScope := c.compileExpr(n.Init)
		c.emitLoc(initIdx, c.rAlloc, initScope)
		reg[0] = c.rAlloc
		c.rAlloc++

		endIdx, endScope := c.compileExpr(n.End)
		c.emitLoc(endIdx, c.rAlloc, endScope)
		reg[1] = c.rAlloc
		c.rAlloc++

		stepIdx, stepScope := c.compileExpr(n.Step)
		c.emitLoc(stepIdx, c.rAlloc, stepScope)
		reg[2] = c.rAlloc
		c.rAlloc++

		stateIdx, stateScope := c.compileExpr(n.State)
		c.sb.addLocal(n.State.(*ast.ForState).Value, c.level, c.scope, c.rAlloc)
		c.emitLoc(stateIdx, c.rAlloc, stateScope)
		reg[3] = c.rAlloc
		c.rAlloc++

		forIndex := c.kb.ForLoopIndex(int(reg[0]), int(reg[1]), int(reg[2]), int(reg[3]))
		c.emitForInit(forIndex, 0)
		jumpAddr := len(c.module.Code)
		postLoopAddr := jumpAddr - 2
		c.compileStmt(n.Block)
		binary.NativeEndian.PutUint16(c.module.Code[postLoopAddr:], uint16(len(c.module.Code)))
		c.emitForLoop(forIndex, jumpAddr)
		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.rAlloc -= 3
		c.scope--
	case *ast.ReferenceStmt:
		idx, scope := c.refScope(n.Value)
		c.emitLoc(idx, c.rAlloc, scope)
	case *ast.IGetStmt:
		resultReg := c.rAlloc
		c.rAlloc++
		j, t := c.compileExpr(n.Index)
		c.rAlloc = resultReg
		c.emitIGet(int(resultReg), j, rLocal, t, resultReg)
	case *ast.SelectStmt:
		resultReg := c.rAlloc
		c.rAlloc++
		j, t := c.compileExpr(n.Selector)
		c.rAlloc = resultReg
		c.emitIGet(int(resultReg), j, rLocal, t, resultReg)
	case *ast.ISet:
		mutableReg := c.rAlloc
		c.rAlloc++
		ii, scopeI := c.compileExpr(n.Index)
		c.rAlloc++
		ie, scopeE := c.compileExpr(n.Expr)
		c.rAlloc = mutableReg
		c.emitISet(ii, ie, scopeI, scopeE, mutableReg, mutableReg)
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
		return c.kb.IntegerIndex(n.Value), rKonst
	case *ast.Float:
		return c.kb.FloatIndex(n.Value), rKonst
	case *ast.String:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.BinaryExpr:
		opReg := c.rAlloc
		c.rAlloc++
		lidx, lscope := c.compileExpr(n.Lhs)
		c.rAlloc++
		ridx, rscope := c.compileExpr(n.Rhs)
		c.rAlloc = opReg
		if lscope == rKonst && rscope == rKonst {
			if val, err := c.kb.Konstants[lidx].Binop(byte(n.Op), c.kb.Konstants[ridx]); err == nil {
				return c.integrateKonst(val)
			}
		}
		switch n.Op {
		case token.EQ, token.NEQ:
			c.emitEq(lidx, ridx, lscope, rscope, opReg, byte(n.Op))
		default:
			c.emitBinary(lidx, ridx, lscope, rscope, opReg, byte(n.Op))
		}
		return int(opReg), rLocal
	case *ast.PrefixExpr:
		idx, scope := c.compileExpr(n.Expr)
		if scope == rKonst {
			if val, err := c.kb.Konstants[idx].Prefix(byte(n.Op)); err == nil {
				return c.integrateKonst(val)
			}
		}
		c.emitPrefix(idx, c.rAlloc, scope, byte(n.Op))
		return int(c.rAlloc), rLocal
	case *ast.Boolean:
		return c.kb.BooleanIndex(n.Value), rKonst
	case *ast.Nil:
		return c.kb.NilIndex(), rKonst
	case *ast.Reference:
		return c.refScope(n.Value)
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
	case *ast.Document:
		if len(n.Pairs) == 0 {
			c.emitDocument(0, c.rAlloc, c.rAlloc)
			return int(c.rAlloc), rLocal
		}
		var count int
		for _, v := range n.Pairs {
			ik, scopeK := c.compileExpr(v.Key)
			c.emitLoc(ik, c.rAlloc, scopeK)
			c.rAlloc++
			iv, scopeV := c.compileExpr(v.Value)
			if scopeV != rLocal {
				c.emitLoc(iv, c.rAlloc, scopeV)
			} else if iv != int(c.rAlloc) {
				c.emitMove(byte(iv), c.rAlloc)
			}
			c.rAlloc++
			count += 2
		}
		c.rAlloc -= byte(count)
		c.emitDocument(byte(count), c.rAlloc, c.rAlloc)
		return int(c.rAlloc), rLocal
	case *ast.Property:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.ForState:
		return c.kb.IntegerIndex(0), rKonst
	case *ast.IGet:
		resultReg := c.rAlloc
		c.rAlloc++
		i, s := c.compileExpr(n.Indexable)
		c.rAlloc++
		j, t := c.compileExpr(n.Index)
		c.rAlloc = resultReg
		c.emitIGet(i, j, s, t, resultReg)
		return int(resultReg), rLocal
	case *ast.Select:
		resultReg := c.rAlloc
		c.rAlloc++
		i, s := c.compileExpr(n.Selectable)
		c.rAlloc++
		j, t := c.compileExpr(n.Selector)
		c.rAlloc = resultReg
		c.emitIGet(i, j, s, t, resultReg)
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

func (c *Compiler) integrateKonst(val Value) (int, byte) {
	switch e := val.(type) {
	case Integer:
		return c.kb.IntegerIndex(int64(e)), rKonst
	case Float:
		return c.kb.FloatIndex(float64(e)), rKonst
	case Bool:
		return c.kb.BooleanIndex(bool(e)), rKonst
	case String:
		return c.kb.StringIndex(e.Value), rKonst
	default:
		return c.kb.NilIndex(), rKonst
	}
}
