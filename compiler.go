package vida

import (
	"encoding/binary"

	"github.com/ever-eduardo/vida/ast"
	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Compiler struct {
	jumps         []int
	breakJumps    []int
	breakCount    []int
	continueJumps []int
	continueCount []int
	fn            []*Function
	ast           *ast.Ast
	module        *Module
	kb            *KonstBuilder
	sb            *symbolBuilder
	scope         int
	level         int
	rAlloc        byte
	hadError      bool
}

func NewCompiler(ast *ast.Ast, moduleName string) *Compiler {
	return &Compiler{
		ast:    ast,
		module: newModule(moduleName),
		kb:     newKonstBuilder(),
		sb:     newSymbolBuilder(),
	}
}

func (c *Compiler) CompileModule() (*Module, error) {
	c.appendHeader()
	for i := range len(c.ast.Statement) {
		c.compileStmt(c.ast.Statement[i])
	}
	if c.hadError {
		return nil, verror.CompilerError
	}
	c.module.Konstants = c.kb.Konstants
	c.appendEnd()
	return c.module, nil
}

func (c *Compiler) compileStmt(node ast.Node) {
	switch n := node.(type) {
	case *ast.Set:
		from, se := c.compileExpr(n.Expr)
		to, si := c.refScope(n.Indentifier)
		if si == rFree {
			c.emitSetF(from, byte(to), se)
		} else if si == rLoc {
			if from == to {
				return
			}
			if se == rLoc {
				c.emitMove(byte(from), byte(to))
			} else {
				c.emitLoc(from, byte(to), se)
			}
		} else if isGlobal := c.sb.isGlobal(n.Indentifier); isGlobal {
			to := c.kb.StringIndex(n.Indentifier)
			c.emitSetG(from, to, se)
		} else {
			c.hadError = true
		}
	case *ast.Let:
		to := c.kb.StringIndex(n.Indentifier)
		c.sb.addGlobal(n.Indentifier)
		from, scope := c.compileExpr(n.Expr)
		c.emitSetG(from, to, scope)
	case *ast.Loc:
		to := c.rAlloc
		from, scope := c.compileExpr(n.Expr)
		c.rAlloc++
		c.sb.addLocal(n.Identifier, c.level, c.scope, to)
		if scope != rLoc {
			c.emitLoc(from, to, scope)
		}
	case *ast.Branch:
		elifCount := len(n.Elifs)
		hasElif := elifCount != 0
		e, hasElse := n.Else.(*ast.Else)
		shouldJumpOutside := hasElif || hasElse
		c.compileConditional(n.If.(*ast.If), shouldJumpOutside)
		if hasElif {
			for i := 0; i < elifCount-1; i++ {
				c.compileConditional(n.Elifs[i].(*ast.If), hasElif)
			}
			c.compileConditional(n.Elifs[elifCount-1].(*ast.If), hasElse)
		}
		if hasElse {
			c.compileStmt(e.Block)
		}
		if shouldJumpOutside {
			addr := len(c.module.Code)
			for _, v := range c.jumps {
				binary.NativeEndian.PutUint16(c.module.Code[v:], uint16(addr))
			}
			c.jumps = c.jumps[:0]
		}
	case *ast.For:
		c.startLoopScope()
		c.scope++

		ireg := c.rAlloc

		initIdx, initScope := c.compileExpr(n.Init)
		c.emitLoc(initIdx, c.rAlloc, initScope)

		c.rAlloc++
		endIdx, endScope := c.compileExpr(n.End)
		c.emitLoc(endIdx, c.rAlloc, endScope)

		c.rAlloc++
		stepIdx, stepScope := c.compileExpr(n.Step)
		c.emitLoc(stepIdx, c.rAlloc, stepScope)

		c.rAlloc++
		c.sb.addLocal(n.Id, c.level, c.scope, c.rAlloc)
		c.emitLoc(c.kb.IntegerIndex(0), c.rAlloc, rKonst)

		c.rAlloc++
		c.emitForSet(ireg, 0)
		jump := len(c.module.Code)
		c.compileStmt(n.Block)
		evalLoopAddr := len(c.module.Code)
		binary.NativeEndian.PutUint16(c.module.Code[jump-2:], uint16(evalLoopAddr))
		c.emitForLoop(ireg, jump)
		c.cleanUpLoopScope(evalLoopAddr)

		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.rAlloc -= 3
		c.scope--
	case *ast.IFor:
		c.startLoopScope()

		c.scope++
		ireg := c.rAlloc
		c.emitLoc(c.kb.IntegerIndex(0), c.rAlloc, rKonst)

		c.rAlloc++
		c.sb.addLocal(n.Key, c.level, c.scope, c.rAlloc)
		c.emitLoc(c.kb.IntegerIndex(0), c.rAlloc, rKonst)

		c.rAlloc++
		c.sb.addLocal(n.Value, c.level, c.scope, c.rAlloc)
		c.emitLoc(c.kb.IntegerIndex(0), c.rAlloc, rKonst)
		c.rAlloc++

		idx, scope := c.compileExpr(n.Expr)
		c.emitIForSet(0, idx, scope, ireg)
		jump := len(c.module.Code)
		c.compileStmt(n.Block)
		evalLoopAddr := len(c.module.Code)
		binary.NativeEndian.PutUint16(c.module.Code[jump-2:], uint16(evalLoopAddr))
		c.emitIForLoop(ireg, jump)
		c.cleanUpLoopScope(evalLoopAddr)

		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.rAlloc--
		c.scope--
	case *ast.While:
		c.startLoopScope()

		init := len(c.module.Code)
		idx, scope := c.compileExpr(n.Condition)
		if scope == rKonst {
			switch v := c.kb.Konstants[idx].(type) {
			case Nil:
				c.skipBlock(n.Block)
				c.cleanUpLoopScope(init)
				return
			case Bool:
				if !v {
					c.skipBlock(n.Block)
					c.cleanUpLoopScope(init)
					return
				}
			}
			c.compileStmt(n.Block)
			c.emitJump(init)
			c.cleanUpLoopScope(init)
		} else {
			addr := len(c.module.Code) + 4
			c.emitTestF(idx, scope, 0)
			c.compileStmt(n.Block)
			c.emitJump(init)
			binary.NativeEndian.PutUint16(c.module.Code[addr:], uint16(len(c.module.Code)))
			c.cleanUpLoopScope(init)
		}
	case *ast.Break:
		c.breakJumps = append(c.breakJumps, len(c.module.Code)+1)
		c.breakCount[len(c.breakCount)-1]++
		c.emitJump(0)
	case *ast.Continue:
		c.continueJumps = append(c.continueJumps, len(c.module.Code)+1)
		c.continueCount[len(c.continueCount)-1]++
		c.emitJump(0)
	case *ast.ReferenceStmt:
		idx, scope := c.refScope(n.Value)
		c.emitLoc(idx, c.rAlloc, scope)
	case *ast.IGetStmt:
		resultReg := c.rAlloc
		c.rAlloc++
		j, t := c.compileExpr(n.Index)
		c.rAlloc = resultReg
		c.emitIGet(int(resultReg), j, rLoc, t, resultReg)
	case *ast.SelectStmt:
		resultReg := c.rAlloc
		c.rAlloc++
		j, t := c.compileExpr(n.Selector)
		c.rAlloc = resultReg
		c.emitIGet(int(resultReg), j, rLoc, t, resultReg)
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
	case *ast.Ret:
		if c.level == 0 {
			c.hadError = true
		}
		i, s := c.compileExpr(n.Expr)
		c.emitRet(i, c.rAlloc, s)
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
			switch n.Op {
			case token.EQ, token.NEQ:
				return c.integrateKonst(c.kb.Konstants[lidx].Equals(c.kb.Konstants[ridx]))
			default:
				if val, err := c.kb.Konstants[lidx].Binop(byte(n.Op), c.kb.Konstants[ridx]); err == nil {
					return c.integrateKonst(val)
				}
			}
		}
		switch n.Op {
		case token.EQ, token.NEQ:
			c.emitEq(lidx, ridx, lscope, rscope, opReg, byte(n.Op))
		default:
			c.emitBinary(lidx, ridx, lscope, rscope, opReg, byte(n.Op))
		}
		return int(opReg), rLoc
	case *ast.PrefixExpr:
		idx, scope := c.compileExpr(n.Expr)
		if scope == rKonst {
			if val, err := c.kb.Konstants[idx].Prefix(byte(n.Op)); err == nil {
				return c.integrateKonst(val)
			}
		}
		c.emitPrefix(idx, c.rAlloc, scope, byte(n.Op))
		return int(c.rAlloc), rLoc
	case *ast.Boolean:
		return c.kb.BooleanIndex(n.Value), rKonst
	case *ast.Nil:
		return c.kb.NilIndex(), rKonst
	case *ast.Reference:
		i, b := c.refScope(n.Value)
		return i, b
	case *ast.List:
		if len(n.ExprList) == 0 {
			c.emitList(0, c.rAlloc, c.rAlloc)
			return int(c.rAlloc), rLoc
		}
		var count int
		for _, v := range n.ExprList {
			idx, scope := c.compileExpr(v)
			if scope != rLoc {
				c.emitLoc(idx, c.rAlloc, scope)
			} else if idx != int(c.rAlloc) {
				c.emitMove(byte(idx), c.rAlloc)
			}
			c.rAlloc++
			count++
		}
		c.rAlloc -= byte(count)
		c.emitList(byte(count), c.rAlloc, c.rAlloc)
		return int(c.rAlloc), rLoc
	case *ast.Document:
		if len(n.Pairs) == 0 {
			c.emitDocument(0, c.rAlloc, c.rAlloc)
			return int(c.rAlloc), rLoc
		}
		var count int
		for _, v := range n.Pairs {
			ik, scopeK := c.compileExpr(v.Key)
			c.emitLoc(ik, c.rAlloc, scopeK)
			c.rAlloc++
			iv, scopeV := c.compileExpr(v.Value)
			if scopeV != rLoc {
				c.emitLoc(iv, c.rAlloc, scopeV)
			} else if iv != int(c.rAlloc) {
				c.emitMove(byte(iv), c.rAlloc)
			}
			c.rAlloc++
			count += 2
		}
		c.rAlloc -= byte(count)
		c.emitDocument(byte(count), c.rAlloc, c.rAlloc)
		return int(c.rAlloc), rLoc
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
		return int(resultReg), rLoc
	case *ast.Select:
		resultReg := c.rAlloc
		c.rAlloc++
		i, s := c.compileExpr(n.Selectable)
		c.rAlloc++
		j, t := c.compileExpr(n.Selector)
		c.rAlloc = resultReg
		c.emitIGet(i, j, s, t, resultReg)
		return int(resultReg), rLoc
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
		return int(resultReg), rLoc
	case *ast.Fun:
		fn := &Function{}
		c.fn = append(c.fn, fn)
		jump := len(c.module.Code)
		c.emitFun(c.kb.FunctionIndex(fn), c.rAlloc, 0)
		fn.First = len(c.module.Code)
		reg := c.startFuncScope()
		for _, v := range n.Args {
			fn.Arity++
			c.sb.addLocal(v, c.level, c.scope, c.rAlloc)
			c.rAlloc++
		}
		c.compileStmt(n.Body)
		binary.NativeEndian.PutUint16(c.module.Code[jump+3:], uint16(len(c.module.Code)))
		c.leaveFuncScope()
		c.rAlloc = reg
		return int(c.rAlloc), rLoc
	default:
		return 0, rKonst
	}
}

func (c *Compiler) compileConditional(n *ast.If, shouldJumpOutside bool) {
	idx, scope := c.compileExpr(n.Condition)
	if scope == rKonst {
		switch v := c.kb.Konstants[idx].(type) {
		case Nil:
			c.skipBlock(n.Block)
			return
		case Bool:
			if !v {
				c.skipBlock(n.Block)
				return
			}
		}
		c.compileBlockAndCheckJump(n.Block, shouldJumpOutside)
	} else {
		addr := len(c.module.Code) + 4
		c.emitTestF(idx, scope, 0)
		c.compileBlockAndCheckJump(n.Block, shouldJumpOutside)
		binary.NativeEndian.PutUint16(c.module.Code[addr:], uint16(len(c.module.Code)))
	}
}

func (c *Compiler) skipBlock(block ast.Node) {
	addr := len(c.module.Code) + 1
	c.emitJump(0)
	c.compileStmt(block)
	binary.NativeEndian.PutUint16(c.module.Code[addr:], uint16(len(c.module.Code)))
}

func (c *Compiler) compileBlockAndCheckJump(block ast.Node, shouldJumpOutside bool) {
	c.compileStmt(block)
	if shouldJumpOutside {
		c.jumps = append(c.jumps, len(c.module.Code)+1)
		c.emitJump(0)
	}
}

func (c *Compiler) cleanUpLoopScope(init int) {
	hasBreaks := len(c.breakJumps)
	lastElem := len(c.breakCount) - 1
	count := c.breakCount[lastElem]
	if hasBreaks > 0 {
		for i := 1; i <= count; i++ {
			binary.NativeEndian.PutUint16(c.module.Code[c.breakJumps[hasBreaks-i]:], uint16(len(c.module.Code)))
		}
		c.breakJumps = c.breakJumps[:hasBreaks-count]
	}
	c.breakCount = c.breakCount[:lastElem]
	hasContinues := len(c.continueJumps)
	lastElem = len(c.continueCount) - 1
	count = c.continueCount[lastElem]
	if hasContinues > 0 {
		for i := 1; i <= count; i++ {
			binary.NativeEndian.PutUint16(c.module.Code[c.continueJumps[hasContinues-i]:], uint16(init))
		}
		c.continueJumps = c.continueJumps[:hasContinues-count]
	}
	c.continueCount = c.continueCount[:lastElem]
}

func (c *Compiler) startLoopScope() {
	c.breakCount = append(c.breakCount, 0)
	c.continueCount = append(c.continueCount, 0)
}

func (c *Compiler) startFuncScope() byte {
	r := c.rAlloc
	c.rAlloc = 0
	c.level++
	return r
}

func (c *Compiler) leaveFuncScope() {
	c.sb.clearLocals(c.level, c.scope)
	c.level--
	c.fn = c.fn[:c.level]
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
