package vida

import (
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
	refStmtIsLoc  [2]int
	fn            []*CoreFunction
	currentFn     *CoreFunction
	ast           *ast.Ast
	module        *Module
	kb            *KonstBuilder
	sb            *symbolBuilder
	scope         int
	level         int
	rAlloc        int
	rDest         int
	mutLoc        bool
	hadError      bool
}

func NewCompiler(ast *ast.Ast, moduleName string) *Compiler {
	c := &Compiler{
		ast:    ast,
		module: newModule(moduleName),
		kb:     newKonstBuilder(),
		sb:     newSymbolBuilder(),
	}
	c.fn = append(c.fn, c.module.MainFunction.CoreFn)
	c.currentFn = c.module.MainFunction.CoreFn
	return c
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
	case *ast.Mut:
		to, sIdent := c.refScope(n.Indentifier)
		switch sIdent {
		case rFree:
			from, sexpr := c.compileExpr(n.Expr)
			switch sexpr {
			case rGlob:
				c.emitLoadG(from, c.rAlloc)
				c.emitStoreF(c.rAlloc, to)
			case rKonst:
				c.emitLoadK(from, c.rAlloc)
				c.emitStoreF(c.rAlloc, to)
			case rFree:
				if from != to {
					c.emitLoadF(from, c.rAlloc)
					c.emitStoreF(c.rAlloc, to)
				}
			case rLoc:
				c.emitStoreF(from, to)
			}
		case rLoc:
			c.mutLoc = true
			c.rDest = to
			from, sexpr := c.compileExpr(n.Expr)
			switch sexpr {
			case rGlob:
				c.emitLoadG(from, to)
			case rLoc:
				if from != to {
					c.emitMove(from, to)
				}
			case rKonst:
				c.emitLoadK(from, to)
			case rFree:
				c.emitLoadF(from, to)
			}
			c.mutLoc = false
		case rGlob:
			from, sexpr := c.compileExpr(n.Expr)
			switch sexpr {
			case rGlob:
				if from != to {
					c.emitLoadG(from, c.rAlloc)
					c.emitStoreG(c.rAlloc, to, 0)
				}
			case rKonst:
				c.emitStoreG(from, to, 1)
			case rFree:
				c.emitLoadF(from, c.rAlloc)
				c.emitStoreG(c.rAlloc, to, 0)
			case rLoc:
				c.emitStoreG(from, to, 0)
			}
		}
	case *ast.Let:
		to := c.sb.addGlobal(n.Indentifier)
		c.module.Store = append(c.module.Store, NilValue)
		from, scope := c.compileExpr(n.Expr)
		switch scope {
		case rKonst:
			c.emitStoreG(from, to, 1)
		case rGlob:
			c.emitLoadG(from, c.rAlloc)
			c.emitStoreG(c.rAlloc, to, 0)
		case rFree:
			c.emitLoadF(from, c.rAlloc)
			c.emitStoreG(c.rAlloc, to, 0)
		case rLoc:
			c.emitStoreG(from, to, 0)
		}
	case *ast.Loc:
		to := c.rAlloc
		from, scope := c.compileExpr(n.Expr)
		c.sb.addLocal(n.Identifier, c.level, c.scope, to)
		switch scope {
		case rKonst:
			c.emitLoadK(from, to)
		case rGlob:
			c.emitLoadG(from, to)
		case rFree:
			c.emitLoadF(from, to)
		case rLoc:
			if from != to {
				c.emitMove(from, to)
			}
		}
		c.rAlloc++
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
			// _ = len(c.currentFn.Code)
			// for _, v := range c.jumps {
			// 	// binary.NativeEndian.PutUint16(c.currentFn.Code[v:], uint16(addr))
			// }
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
		jump := len(c.currentFn.Code)
		c.compileStmt(n.Block)
		// evalLoop_ = len(c.currentFn.Code)
		// binary.NativeEndian.PutUint16(c.currentFn.Code[jump-2:], uint16(evalLoopAddr))
		c.emitForLoop(ireg, jump)
		// c.cleanUpLoopScope(evalLoopAddr)

		c.rAlloc -= c.sb.clearLocals(c.level, c.scope)
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
		jump := len(c.currentFn.Code)
		c.compileStmt(n.Block)
		// evalLoop_ = len(c.currentFn.Code)
		// binary.NativeEndian.PutUint16(c.currentFn.Code[jump-2:], uint16(evalLoopAddr))
		c.emitIForLoop(ireg, jump)
		// c.cleanUpLoopScope(evalLoopAddr)

		c.rAlloc -= c.sb.clearLocals(c.level, c.scope)
		c.rAlloc--
		c.scope--
	case *ast.While:
		c.startLoopScope()

		init := len(c.currentFn.Code)
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
			// _ = len(c.currentFn.Code) + 4
			c.emitTestF(idx, scope, 0)
			c.compileStmt(n.Block)
			c.emitJump(init)
			// binary.NativeEndian.PutUint16(c.currentFn.Code[addr:], uint16(len(c.currentFn.Code)))
			c.cleanUpLoopScope(init)
		}
	case *ast.Break:
		c.breakJumps = append(c.breakJumps, len(c.currentFn.Code)+1)
		c.breakCount[len(c.breakCount)-1]++
		c.emitJump(0)
	case *ast.Continue:
		c.continueJumps = append(c.continueJumps, len(c.currentFn.Code)+1)
		c.continueCount[len(c.continueCount)-1]++
		c.emitJump(0)
	case *ast.ReferenceStmt:
		i, s := c.refScope(n.Value)
		switch s {
		case rLoc:
			c.refStmtIsLoc[0], c.refStmtIsLoc[1] = 1, i
		case rGlob:
			c.emitLoadG(i, c.rAlloc)
		case rFree:
			c.emitLoadF(i, c.rAlloc)
		}
	case *ast.IGetStmt:
		if c.refStmtIsLoc[0] == 1 {
			c.refStmtIsLoc[0] = 0
			i := c.refStmtIsLoc[1]
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				c.emitIGet(i, j, c.rAlloc, 0)
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
			case rKonst:
				c.emitIGet(i, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
			}
		} else {
			j, t := c.compileExpr(n.Index)
			switch t {
			case rKonst:
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rLoc:
				c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
			case rGlob:
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rFree:
				c.emitLoadF(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		}
	case *ast.SelectStmt:
		if c.refStmtIsLoc[0] == 1 {
			c.refStmtIsLoc[0] = 0
			i := c.refStmtIsLoc[1]
			j, t := c.compileExpr(n.Selector)
			switch t {
			case rLoc:
				c.emitIGet(i, j, c.rAlloc, 0)
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
			case rKonst:
				c.emitIGet(i, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
			}
		} else {
			j, t := c.compileExpr(n.Selector)
			switch t {
			case rKonst:
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rLoc:
				c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
			case rGlob:
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rFree:
				c.emitLoadF(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		}
	case *ast.ISet:
		if c.refStmtIsLoc[0] == 1 {
			c.refStmtIsLoc[0] = 0
			i := c.refStmtIsLoc[1]
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				c.rAlloc++
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitISet(i, j, k, 0)
				case rKonst:
					c.emitISet(i, j, k, 1)
				case rGlob:
					c.emitLoadG(k, c.rAlloc)
					c.emitISet(i, j, c.rAlloc, 0)
				case rFree:
					c.emitLoadF(k, c.rAlloc)
					c.emitISet(i, j, c.rAlloc, 0)
				}
				c.rAlloc--
			case rGlob:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitLoadG(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 0)
				case rKonst:
					c.emitLoadG(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 1)
				case rGlob:
					c.emitLoadG(j, c.rAlloc)
					c.emitLoadG(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				case rFree:
					c.emitLoadG(j, c.rAlloc)
					c.emitLoadF(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				}
			case rKonst:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitISetK(i, j, k, 0)
				case rKonst:
					c.emitISetK(i, j, k, 1)
				case rGlob:
					c.emitLoadG(k, c.rAlloc)
					c.emitISetK(i, j, c.rAlloc, 0)
				case rFree:
					c.emitLoadF(k, c.rAlloc)
					c.emitISetK(i, j, c.rAlloc, 0)
				}
			case rFree:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitLoadF(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 0)
				case rKonst:
					c.emitLoadF(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 1)
				case rGlob:
					c.emitLoadF(j, c.rAlloc)
					c.emitLoadG(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				case rFree:
					c.emitLoadF(j, c.rAlloc)
					c.emitLoadF(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				}
			}
		} else {
			i := c.rAlloc
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				c.rAlloc++
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitISet(i, j, k, 0)
				case rKonst:
					c.emitISet(i, j, k, 1)
				case rGlob:
					c.emitLoadG(k, c.rAlloc)
					c.emitISet(i, j, c.rAlloc, 0)
				case rFree:
					c.emitLoadF(k, c.rAlloc)
					c.emitISet(i, j, c.rAlloc, 0)
				}
				c.rAlloc--
			case rGlob:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitLoadG(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 0)
				case rKonst:
					c.emitLoadG(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 1)
				case rGlob:
					c.emitLoadG(j, c.rAlloc)
					c.emitLoadG(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				case rFree:
					c.emitLoadG(j, c.rAlloc)
					c.emitLoadF(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				}
			case rKonst:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitISetK(i, j, k, 0)
				case rKonst:
					c.emitISetK(i, j, k, 1)
				case rGlob:
					c.emitLoadG(k, c.rAlloc)
					c.emitISetK(i, j, c.rAlloc, 0)
				case rFree:
					c.emitLoadF(k, c.rAlloc)
					c.emitISetK(i, j, c.rAlloc, 0)
				}
			case rFree:
				k, u := c.compileExpr(n.Expr)
				switch u {
				case rLoc:
					c.emitLoadF(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 0)
				case rKonst:
					c.emitLoadF(j, c.rAlloc)
					c.emitISet(i, c.rAlloc, k, 1)
				case rGlob:
					c.emitLoadF(j, c.rAlloc)
					c.emitLoadG(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				case rFree:
					c.emitLoadF(j, c.rAlloc)
					c.emitLoadF(k, c.rAlloc+1)
					c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
				}
			}
		}
	case *ast.Block:
		c.scope++
		for i := range len(n.Statement) {
			c.compileStmt(n.Statement[i])
		}
		locals := c.sb.clearLocals(c.level, c.scope)
		c.rAlloc -= locals
		c.scope--
	case *ast.Ret:
		if c.level == 0 {
			c.hadError = true
		}
		i, s := c.compileExpr(n.Expr)
		c.emitRet(i, s)
	case *ast.CallStmt:
		reg := c.rAlloc
		for _, v := range n.Args {
			c.rAlloc++
			i, s := c.compileExpr(v)
			c.emitLoc(i, c.rAlloc, s)
		}
		c.rAlloc = reg
		c.emitCall(reg, len(n.Args))
	case *ast.MethodCallStmt:
		// reg := c.rAlloc
		// c.rAlloc++
		// c.emitMove(reg, c.rAlloc)
		// c.rAlloc++
		// j, t := c.compileExpr(n.Prop)
		// c.emitIGet(reg, j, rLoc, t, reg)
		// for _, v := range n.Args {
		// 	i, s := c.compileExpr(v)
		// 	c.emitLoc(i, c.rAlloc, s)
		// 	c.rAlloc++
		// }
		// c.rAlloc = reg
		// c.emitCall(reg, len(n.Args)+1)
	}
}

func (c *Compiler) compileExpr(node ast.Node) (int, int) {
	switch n := node.(type) {
	case *ast.Integer:
		return c.kb.IntegerIndex(n.Value), rKonst
	case *ast.Float:
		return c.kb.FloatIndex(n.Value), rKonst
	case *ast.String:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.BinaryExpr:
		lidx, lscope := c.compileExpr(n.Lhs)
		lreg := c.rAlloc
		switch lscope {
		case rKonst:
			ridx, rscope := c.compileExpr(n.Rhs)
			switch rscope {
			case rKonst:
				switch n.Op {
				case token.EQ, token.NEQ:
					return c.integrateKonst(c.kb.Konstants[lidx].Equals(c.kb.Konstants[ridx]))
				default:
					if val, err := c.kb.Konstants[lidx].Binop(uint64(n.Op), c.kb.Konstants[ridx]); err == nil {
						return c.integrateKonst(val)
					} else {
						c.hadError = true
					}
				}
			case rGlob:
				c.emitLoadG(ridx, lreg)
				c.emitBinopK(lidx, lreg, lreg, n.Op)
			case rLoc:
				c.emitBinopK(lidx, ridx, lreg, n.Op)
			case rFree:
				c.emitLoadF(ridx, lreg)
				c.emitBinopK(lidx, lreg, lreg, n.Op)
			}
		case rGlob:
			ridx, rscope := c.compileExpr(n.Rhs)
			switch rscope {
			case rGlob:
				c.emitBinopG(lidx, ridx, lreg, n.Op)
			case rKonst:
				c.emitLoadG(lidx, lreg)
				c.emitBinopK(ridx, lreg, lreg, n.Op)
			case rLoc:
				if c.mutLoc && c.rDest == ridx {
					c.emitLoadG(lidx, lreg)
					c.emitBinop(lreg, ridx, c.rDest, n.Op)
					return c.rDest, rLoc
				} else {
					c.emitLoadG(lidx, lreg)
					c.emitBinop(lreg, ridx, lreg, n.Op)
				}
			case rFree:
				c.rAlloc++
				c.emitLoadG(lidx, lreg)
				c.emitLoadF(ridx, c.rAlloc)
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		case rLoc:
			c.rAlloc++
			ridx, rscope := c.compileExpr(n.Rhs)
			switch rscope {
			case rLoc:
				if c.mutLoc && (c.rDest == lidx || c.rDest == ridx) {
					c.emitBinop(lidx, ridx, c.rDest, n.Op)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitBinop(lidx, ridx, lreg, n.Op)
				}
			case rGlob:
				if c.mutLoc && (c.rDest == lidx) {
					c.emitLoadG(ridx, c.rAlloc)
					c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitLoadG(ridx, c.rAlloc)
					c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				}
			case rKonst:
				if c.mutLoc && (c.rDest == lidx) {
					c.emitBinopK(ridx, lidx, c.rDest, n.Op)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitBinopK(lidx, ridx, lreg, n.Op)
				}
			case rFree:
				if c.mutLoc && (c.rDest == lidx) {
					c.emitLoadF(ridx, c.rAlloc)
					c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitLoadF(ridx, c.rAlloc)
					c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				}
			}
			c.rAlloc--
		case rFree:
			ridx, rscope := c.compileExpr(n.Rhs)
			switch rscope {
			case rLoc:
				if c.mutLoc && (c.rDest == ridx) {
					c.emitLoadF(lidx, lreg)
					c.emitBinop(lreg, ridx, c.rDest, n.Op)
					return c.rDest, rLoc
				} else {
					c.emitLoadF(lidx, lreg)
					c.emitBinop(lreg, ridx, lreg, n.Op)
				}
			case rGlob:
				c.rAlloc++
				c.emitLoadF(lidx, lreg)
				c.emitLoadG(ridx, c.rAlloc)
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			case rKonst:
				c.emitLoadF(lidx, lreg)
				c.emitBinopK(ridx, lreg, lreg, n.Op)
			case rFree:
				c.rAlloc++
				c.emitLoadF(lidx, lreg)
				c.emitLoadF(ridx, c.rAlloc)
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
		return lreg, rLoc
	case *ast.PrefixExpr:
		from, scope := c.compileExpr(n.Expr)
		if scope == rKonst {
			if val, err := c.kb.Konstants[from].Prefix(uint64(n.Op)); err == nil {
				return c.integrateKonst(val)
			}
		}
		switch scope {
		case rGlob:
			c.emitLoadG(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
		case rLoc:
			if c.mutLoc && c.rDest == from {
				c.emitPrefix(from, from, n.Op)
				return from, rLoc
			} else {
				c.emitPrefix(from, c.rAlloc, n.Op)
			}
		case rKonst:
			c.emitLoadK(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
		case rFree:
			c.emitLoadF(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
		}
		return c.rAlloc, rLoc
	case *ast.Boolean:
		return c.kb.BooleanIndex(n.Value), rKonst
	case *ast.Nil:
		return c.kb.NilIndex(), rKonst
	case *ast.Reference:
		return c.refScope(n.Value)
	case *ast.List:
		if len(n.ExprList) == 0 {
			c.emitList(0, c.rAlloc)
			return c.rAlloc, rLoc
		}
		var count int
		for _, v := range n.ExprList {
			i, s := c.compileExpr(v)
			switch s {
			case rLoc:
				if i != c.rAlloc {
					c.emitMove(i, c.rAlloc)
				}
			case rKonst:
				c.emitLoadK(i, c.rAlloc)
			case rGlob:
				c.emitLoadG(i, c.rAlloc)
			case rFree:
				c.emitLoadF(i, c.rAlloc)
			}
			c.rAlloc++
			count++
		}
		c.rAlloc -= count
		c.emitList(count, c.rAlloc)
		return c.rAlloc, rLoc
	case *ast.Object:
		c.emitObject(c.rAlloc)
		for _, v := range n.Pairs {
			k, _ := c.compileExpr(v.Key)
			v, sv := c.compileExpr(v.Value)
			switch sv {
			case rKonst:
				c.emitISetK(c.rAlloc, k, v, 1)
			case rLoc:
				c.emitISetK(c.rAlloc, k, v, 0)
			case rGlob:
				c.rAlloc++
				c.emitLoadG(v, c.rAlloc)
				c.emitISetK(c.rAlloc-1, k, c.rAlloc, 0)
				c.rAlloc--
			case rFree:
				c.rAlloc++
				c.emitLoadF(v, c.rAlloc)
				c.emitISetK(c.rAlloc-1, k, c.rAlloc, 0)
				c.rAlloc--
			}
		}
		return c.rAlloc, rLoc
	case *ast.Property:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.ForState:
		return c.kb.IntegerIndex(0), rKonst
	case *ast.IGet:
		i, s := c.compileExpr(n.Indexable)
		switch s {
		case rLoc:
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitIGet(i, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, c.rAlloc, 0)
				}
			case rGlob:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitLoadG(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadG(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
				}
			case rKonst:
				if c.mutLoc && (i == c.rDest) {
					c.emitIGet(i, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, c.rAlloc, 1)
				}
			case rFree:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitLoadF(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadF(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
				}
			}
		case rGlob:
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				if c.mutLoc && j == c.rDest {
					c.emitLoadG(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadG(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadG(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rKonst:
				c.emitLoadG(i, c.rAlloc)
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadG(i, c.rAlloc)
				c.emitLoadF(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		case rFree:
			j, t := c.compileExpr(n.Index)
			switch t {
			case rLoc:
				if c.mutLoc && j == c.rDest {
					c.emitLoadF(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadF(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadF(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rKonst:
				c.emitLoadG(i, c.rAlloc)
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadF(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		}
		return c.rAlloc, rLoc
	case *ast.Select:
		i, s := c.compileExpr(n.Selectable)
		switch s {
		case rLoc:
			j, t := c.compileExpr(n.Selector)
			switch t {
			case rLoc:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitIGet(i, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, c.rAlloc, 0)
				}
			case rGlob:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitLoadG(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadG(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
				}
			case rKonst:
				if c.mutLoc && (i == c.rDest) {
					c.emitIGet(i, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, c.rAlloc, 1)
				}
			case rFree:
				if c.mutLoc && (i == c.rDest || j == c.rDest) {
					c.emitLoadF(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadF(j, c.rAlloc)
					c.emitIGet(i, c.rAlloc, c.rAlloc, 0)
				}
			}
		case rGlob:
			j, t := c.compileExpr(n.Selector)
			switch t {
			case rLoc:
				if c.mutLoc && j == c.rDest {
					c.emitLoadG(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadG(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadG(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rKonst:
				c.emitLoadG(i, c.rAlloc)
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadG(i, c.rAlloc)
				c.emitLoadF(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		case rFree:
			j, t := c.compileExpr(n.Selector)
			switch t {
			case rLoc:
				if c.mutLoc && j == c.rDest {
					c.emitLoadF(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitLoadF(i, c.rAlloc)
					c.emitIGet(c.rAlloc, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadF(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			case rKonst:
				c.emitLoadG(i, c.rAlloc)
				c.emitIGet(c.rAlloc, j, c.rAlloc, 1)
			case rFree:
				c.emitLoadF(i, c.rAlloc)
				c.emitLoadG(j, c.rAlloc+1)
				c.emitIGet(c.rAlloc, c.rAlloc+1, c.rAlloc, 0)
			}
		}
		return c.rAlloc, rLoc
	case *ast.Slice:
		resultReg := c.rAlloc
		var scopeV, scopeL, scopeR int
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
		c.emitSlice(n.Mode, fromV, fromL, fromR, scopeV, scopeL, scopeR, resultReg)
		return resultReg, rLoc
	case *ast.Fun:
		fn := &CoreFunction{}
		c.fn = append(c.fn, fn)
		c.emitFun(c.kb.FunctionIndex(fn), c.rAlloc)
		c.currentFn = fn
		reg := c.startFuncScope()
		for _, v := range n.Args {
			fn.Arity++
			c.sb.addLocal(v, c.level, c.scope, c.rAlloc)
			c.rAlloc++
		}
		c.compileStmt(n.Body)
		c.leaveFuncScope()
		c.rAlloc = reg
		return c.rAlloc, rLoc
	case *ast.CallExpr:
		reg := c.rAlloc
		idx, s := c.compileExpr(n.Fun)
		c.emitLoc(idx, reg, s)
		for _, v := range n.Args {
			c.rAlloc++
			i, s := c.compileExpr(v)
			c.emitLoc(i, c.rAlloc, s)
		}
		c.rAlloc = reg
		c.emitCall(reg, len(n.Args))
		return reg, rLoc
	case *ast.MethodCallExpr:
		// reg := c.rAlloc
		// c.rAlloc++
		// i, s := c.compileExpr(n.Obj)
		// c.rAlloc++
		// j, t := c.compileExpr(n.Prop)
		// c.rAlloc = reg
		// c.emitIGet(i, j, s, t, reg)
		// c.rAlloc++
		// i, s = c.compileExpr(n.Obj)
		// c.emitLoc(i, c.rAlloc, s)
		// c.rAlloc++
		// for _, v := range n.Args {
		// 	i, s := c.compileExpr(v)
		// 	c.emitLoc(i, c.rAlloc, s)
		// 	c.rAlloc++
		// }
		// c.rAlloc = reg
		// c.emitCall(reg, len(n.Args)+1)
		return 0, rLoc
	default:
		return 0, rGlob
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
		_ = len(c.currentFn.Code) + 4
		c.emitTestF(idx, scope, 0)
		c.compileBlockAndCheckJump(n.Block, shouldJumpOutside)
		// binary.NativeEndian.PutUint16(c.currentFn.Code[addr:], uint16(len(c.currentFn.Code)))
	}
}

func (c *Compiler) skipBlock(block ast.Node) {
	_ = len(c.currentFn.Code) + 1
	c.emitJump(0)
	c.compileStmt(block)
	// binary.NativeEndian.PutUint16(c.currentFn.Code[addr:], uint16(len(c.currentFn.Code)))
}

func (c *Compiler) compileBlockAndCheckJump(block ast.Node, shouldJumpOutside bool) {
	c.compileStmt(block)
	if shouldJumpOutside {
		c.jumps = append(c.jumps, len(c.currentFn.Code)+1)
		c.emitJump(0)
	}
}

func (c *Compiler) cleanUpLoopScope(init int) {
	hasBreaks := len(c.breakJumps)
	lastElem := len(c.breakCount) - 1
	count := c.breakCount[lastElem]
	if hasBreaks > 0 {
		for i := 1; i <= count; i++ {
			// binary.NativeEndian.PutUint16(c.currentFn.Code[c.breakJumps[hasBreaks-i]:], uint16(len(c.currentFn.Code)))
		}
		c.breakJumps = c.breakJumps[:hasBreaks-count]
	}
	c.breakCount = c.breakCount[:lastElem]
	hasContinues := len(c.continueJumps)
	lastElem = len(c.continueCount) - 1
	count = c.continueCount[lastElem]
	if hasContinues > 0 {
		for i := 1; i <= count; i++ {
			// binary.NativeEndian.PutUint16(c.currentFn.Code[c.continueJumps[hasContinues-i]:], uint16(init))
		}
		c.continueJumps = c.continueJumps[:hasContinues-count]
	}
	c.continueCount = c.continueCount[:lastElem]
}

func (c *Compiler) startLoopScope() {
	c.breakCount = append(c.breakCount, 0)
	c.continueCount = append(c.continueCount, 0)
}

func (c *Compiler) startFuncScope() int {
	r := c.rAlloc
	c.rAlloc = 0
	c.level++
	return r
}

func (c *Compiler) leaveFuncScope() {
	c.sb.clearLocals(c.level, c.scope)
	c.fn = c.fn[:c.level]
	c.level--
	c.currentFn = c.fn[c.level]
}

func (c *Compiler) integrateKonst(val Value) (int, int) {
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
