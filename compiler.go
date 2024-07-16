package vida

import (
	"encoding/binary"

	"github.com/ever-eduardo/vida/ast"
	"github.com/ever-eduardo/vida/token"
)

type Compiler struct {
	jumps         []int
	breakJumps    []int
	breakCount    []int
	continueJumps []int
	continueCount []int
	ast           *ast.Ast
	module        *Module
	function      *Function
	parent        *Compiler
	kb            *KonstBuilder
	sb            *symbolBuilder
	scope         int
	level         int
	rAlloc        byte
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
		c.breakCount = append(c.breakCount, 0)
		c.continueCount = append(c.continueCount, 0)
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
		c.emitForSet(forIndex, 0)
		jump := len(c.module.Code)
		c.compileStmt(n.Block)
		evalLoopAddr := len(c.module.Code)
		binary.NativeEndian.PutUint16(c.module.Code[jump-2:], uint16(evalLoopAddr))
		c.emitForLoop(forIndex, jump)
		c.cleanUpJumps(evalLoopAddr)

		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.rAlloc -= 3
		c.scope--
	case *ast.IFor:
		c.breakCount = append(c.breakCount, 0)
		c.continueCount = append(c.continueCount, 0)
		c.scope++
		var reg [4]byte

		initIdx, initScope := c.compileExpr(&ast.ForState{Value: iter})
		c.emitLoc(initIdx, c.rAlloc, initScope)
		reg[0] = c.rAlloc
		c.rAlloc++

		endIdx, endScope := c.compileExpr(&ast.ForState{Value: state})
		c.emitLoc(endIdx, c.rAlloc, endScope)
		reg[1] = c.rAlloc
		c.rAlloc++

		stepIdx, stepScope := c.compileExpr(n.Key)
		c.sb.addLocal(n.Key.(*ast.ForState).Value, c.level, c.scope, c.rAlloc)
		c.emitLoc(stepIdx, c.rAlloc, stepScope)
		reg[2] = c.rAlloc
		c.rAlloc++

		stateIdx, stateScope := c.compileExpr(n.Value)
		c.sb.addLocal(n.Value.(*ast.ForState).Value, c.level, c.scope, c.rAlloc)
		c.emitLoc(stateIdx, c.rAlloc, stateScope)
		reg[3] = c.rAlloc
		c.rAlloc++

		forIndex := c.kb.IForLoopIndex(int(reg[0]), int(reg[1]), int(reg[2]), int(reg[3]))
		idx, scope := c.compileExpr(n.Expr)
		c.emitIForSet(0, idx, scope, reg[0])
		jump := len(c.module.Code)
		c.compileStmt(n.Block)
		evalLoopAddr := len(c.module.Code)
		binary.NativeEndian.PutUint16(c.module.Code[jump-2:], uint16(evalLoopAddr))
		c.emitIForLoop(forIndex, jump)
		c.cleanUpJumps(evalLoopAddr)

		c.rAlloc -= byte(c.sb.clearLocals(c.level, c.scope))
		c.rAlloc -= 2
		c.scope--
	case *ast.While:
		c.breakCount = append(c.breakCount, 0)
		c.continueCount = append(c.continueCount, 0)
		init := len(c.module.Code)
		idx, scope := c.compileExpr(n.Condition)
		if scope == rKonst {
			switch v := c.kb.Konstants[idx].(type) {
			case Nil:
				c.skipBlock(n.Block)
				c.cleanUpJumps(init)
				return
			case Bool:
				if !v {
					c.skipBlock(n.Block)
					c.cleanUpJumps(init)
					return
				}
			}
			c.compileStmt(n.Block)
			c.emitJump(init)
			c.cleanUpJumps(init)
		} else {
			addr := len(c.module.Code) + 4
			c.emitTestF(idx, scope, 0)
			c.compileStmt(n.Block)
			c.emitJump(init)
			binary.NativeEndian.PutUint16(c.module.Code[addr:], uint16(len(c.module.Code)))
			c.cleanUpJumps(init)
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

func (c *Compiler) cleanUpJumps(init int) {
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
