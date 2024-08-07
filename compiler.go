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
	fromRefStmt   bool
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
			from, sexpr := c.compileExpr(n.Expr, true)
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
			from, sexpr := c.compileExpr(n.Expr, true)
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
			from, sexpr := c.compileExpr(n.Expr, true)
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
		to, isPresent := c.sb.addGlobal(n.Indentifier)
		if !isPresent {
			c.module.Store = append(c.module.Store, NilValue)
		}
		from, scope := c.compileExpr(n.Expr, true)
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
		from, scope := c.compileExpr(n.Expr, true)
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
			addr := len(c.currentFn.Code)
			for _, v := range c.jumps {
				c.currentFn.Code[v] |= uint64(addr)
			}
			c.jumps = c.jumps[:0]
		}
	case *ast.For:
		c.startLoopScope()
		c.scope++
		ireg := c.rAlloc

		initIdx, initScope := c.compileExpr(n.Init, true)
		c.exprToReg(initIdx, initScope)

		c.rAlloc++
		endIdx, endScope := c.compileExpr(n.End, true)
		c.exprToReg(endIdx, endScope)

		c.rAlloc++
		stepIdx, stepScope := c.compileExpr(n.Step, true)
		c.exprToReg(stepIdx, stepScope)

		c.rAlloc++
		c.sb.addLocal(n.Id, c.level, c.scope, c.rAlloc)
		c.emitLoadK(c.kb.IntegerIndex(0), c.rAlloc)

		c.rAlloc++
		c.emitForSet(ireg, 0)
		loop := len(c.currentFn.Code)

		c.compileStmt(n.Block)
		checkLoop := len(c.currentFn.Code)

		c.currentFn.Code[loop-1] |= uint64(checkLoop) << shift16
		c.emitForLoop(ireg, loop)
		c.cleanUpLoopScope(loop, false)

		c.rAlloc -= c.sb.clearLocals(c.level, c.scope)
		c.rAlloc -= 3
		c.scope--
	case *ast.IFor:
		c.startLoopScope()
		c.scope++
		ireg := c.rAlloc
		c.emitLoadK(c.kb.IntegerIndex(0), ireg)

		c.rAlloc++
		c.sb.addLocal(n.Key, c.level, c.scope, c.rAlloc)
		c.emitLoadK(c.kb.IntegerIndex(0), c.rAlloc)

		c.rAlloc++
		c.sb.addLocal(n.Value, c.level, c.scope, c.rAlloc)
		c.emitLoadK(c.kb.IntegerIndex(0), c.rAlloc)

		c.rAlloc++
		i, s := c.compileExpr(n.Expr, true)
		c.exprToReg(i, s)

		c.emitIForSet(0, c.rAlloc, ireg)
		loop := len(c.currentFn.Code)

		c.compileStmt(n.Block)
		checkLoop := len(c.currentFn.Code)

		c.currentFn.Code[loop-1] |= uint64(checkLoop) << shift32
		c.emitIForLoop(ireg, loop)
		c.cleanUpLoopScope(loop, false)

		c.rAlloc -= c.sb.clearLocals(c.level, c.scope)
		c.rAlloc--
		c.scope--
	case *ast.While:
		c.startLoopScope()
		init := len(c.currentFn.Code)
		idx, scope := c.compileExpr(n.Condition, true)
		if scope == rKonst {
			switch v := c.kb.Konstants[idx].(type) {
			case Nil:
				c.skipBlock(n.Block)
				c.cleanUpLoopScope(init, true)
				return
			case Bool:
				if !v {
					c.skipBlock(n.Block)
					c.cleanUpLoopScope(init, true)
					return
				}
			}
			c.compileStmt(n.Block)
			c.emitJump(init)
			c.cleanUpLoopScope(init, true)
		} else {
			addr := len(c.currentFn.Code)
			c.emitCheck(0, idx, 0)
			c.compileStmt(n.Block)
			c.emitJump(init)
			c.currentFn.Code[addr] |= uint64(len(c.currentFn.Code))
			c.cleanUpLoopScope(init, true)
		}
	case *ast.Break:
		c.breakJumps = append(c.breakJumps, len(c.currentFn.Code))
		c.breakCount[len(c.breakCount)-1]++
		c.emitJump(0)
	case *ast.Continue:
		c.continueJumps = append(c.continueJumps, len(c.currentFn.Code))
		c.continueCount[len(c.continueCount)-1]++
		c.emitJump(0)
	case *ast.ReferenceStmt:
		c.fromRefStmt = true
		i, s := c.refScope(n.Value)
		switch s {
		case rLoc:
			c.emitMove(i, c.rAlloc)
		case rGlob:
			c.emitLoadG(i, c.rAlloc)
		case rFree:
			c.emitLoadF(i, c.rAlloc)
		}
		c.rAlloc++
	case *ast.IGetStmt:
		i := c.rAlloc
		if c.fromRefStmt {
			i -= 1
		} else {
			c.rAlloc++
		}
		j, t := c.compileExpr(n.Index, true)
		switch t {
		case rKonst:
			c.emitIGet(i, j, i, 1)
		case rLoc:
			c.emitIGet(i, j, i, 0)
		case rGlob:
			c.emitLoadG(j, c.rAlloc)
			c.emitIGet(i, c.rAlloc, i, 0)
		case rFree:
			c.emitLoadF(j, c.rAlloc)
			c.emitIGet(i, c.rAlloc, i, 0)
		}
		if !c.fromRefStmt {
			c.rAlloc--
		}
	case *ast.SelectStmt:
		i := c.rAlloc
		if c.fromRefStmt {
			i -= 1
		} else {
			c.rAlloc++
		}
		j, t := c.compileExpr(n.Selector, true)
		switch t {
		case rKonst:
			c.emitIGet(i, j, i, 1)
		case rLoc:
			c.emitIGet(i, j, i, 0)
		case rGlob:
			c.emitLoadG(j, c.rAlloc)
			c.emitIGet(i, c.rAlloc, i, 0)
		case rFree:
			c.emitLoadF(j, c.rAlloc)
			c.emitIGet(i, c.rAlloc, i, 0)
		}
		if !c.fromRefStmt {
			c.rAlloc--
		}
	case *ast.ISet:
		i := c.rAlloc
		if c.fromRefStmt {
			i -= 1
		} else {
			c.rAlloc++
		}
		j, t := c.compileExpr(n.Index, true)
		switch t {
		case rLoc:
			c.rAlloc++
			k, u := c.compileExpr(n.Expr, true)
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
			c.rAlloc++
			k, u := c.compileExpr(n.Expr, true)
			c.rAlloc++
			c.emitLoadG(j, c.rAlloc)
			switch u {
			case rLoc:
				c.emitISet(i, c.rAlloc, k, 0)
			case rKonst:
				c.emitISet(i, c.rAlloc, k, 1)
			case rGlob:
				c.emitLoadG(k, c.rAlloc+1)
				c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
			case rFree:
				c.emitLoadF(k, c.rAlloc+1)
				c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
			}
			c.rAlloc -= 2
		case rKonst:
			c.rAlloc++
			k, u := c.compileExpr(n.Expr, true)
			switch u {
			case rLoc:
				c.emitISetK(i, j, k, 0)
			case rKonst:
				c.emitISetK(i, j, k, 1)
			case rGlob:
				c.rAlloc++
				c.emitLoadG(k, c.rAlloc)
				c.emitISetK(i, j, c.rAlloc, 0)
				c.rAlloc--
			case rFree:
				c.rAlloc++
				c.emitLoadF(k, c.rAlloc)
				c.emitISetK(i, j, c.rAlloc, 0)
				c.rAlloc--
			}
			c.rAlloc--
		case rFree:
			c.rAlloc++
			k, u := c.compileExpr(n.Expr, true)
			c.rAlloc++
			c.emitLoadF(j, c.rAlloc)
			switch u {
			case rLoc:
				c.emitISet(i, c.rAlloc, k, 0)
			case rKonst:
				c.emitISet(i, c.rAlloc, k, 1)
			case rGlob:
				c.emitLoadG(k, c.rAlloc+1)
				c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
			case rFree:
				c.emitLoadF(k, c.rAlloc+1)
				c.emitISet(i, c.rAlloc, c.rAlloc+1, 0)
			}
			c.rAlloc -= 2
		}
		c.rAlloc--
		c.fromRefStmt = false
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
		i, s := c.compileExpr(n.Expr, true)
		c.exprToReg(i, s)
		c.emitRet(c.rAlloc)
	case *ast.CallStmt:
		callable := c.rAlloc
		if c.fromRefStmt {
			callable -= 1
		} else {
			c.rAlloc++
		}
		for _, v := range n.Args {
			i, s := c.compileExpr(v, true)
			c.exprToReg(i, s)
			c.rAlloc++
		}
		c.rAlloc = callable
		c.fromRefStmt = false
		c.emitCall(callable, len(n.Args))
	case *ast.MethodCallStmt:
		obj := c.rAlloc
		if c.fromRefStmt {
			obj -= 1
		} else {
			c.rAlloc++
		}
		c.emitMove(obj, c.rAlloc)
		c.rAlloc++
		j, t := c.compileExpr(n.Prop, true)
		c.exprToReg(j, t)
		c.emitIGet(obj+1, obj+2, obj, 0)
		for _, v := range n.Args {
			i, s := c.compileExpr(v, true)
			c.exprToReg(i, s)
			c.rAlloc++
		}
		c.rAlloc = obj
		c.fromRefStmt = false
		c.emitCall(obj, len(n.Args)+1)
	}
}

func (c *Compiler) compileExpr(node ast.Node, isRoot bool) (int, int) {
	switch n := node.(type) {
	case *ast.Integer:
		return c.kb.IntegerIndex(n.Value), rKonst
	case *ast.Float:
		return c.kb.FloatIndex(n.Value), rKonst
	case *ast.String:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.BinaryExpr:
		switch n.Op {
		case token.EQ, token.NEQ:
			return c.compileBinaryEq(n, isRoot)
		default:
			return c.compileBinaryExpr(n, isRoot)
		}
	case *ast.PrefixExpr:
		from, scope := c.compileExpr(n.Expr, false)
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
			if c.mutLoc && isRoot {
				c.emitPrefix(from, c.rDest, n.Op)
				return c.rDest, rLoc
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
		var count int
		for _, v := range n.ExprList {
			i, s := c.compileExpr(v, false)
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
		if c.mutLoc && isRoot {
			c.emitList(count, c.rAlloc, c.rDest)
			return c.rDest, rLoc
		}
		c.emitList(count, c.rAlloc, c.rAlloc)
		return c.rAlloc, rLoc
	case *ast.Object:
		objAddr := c.rAlloc
		if c.mutLoc && isRoot {
			objAddr = c.rDest
		}
		c.emitObject(objAddr)
		for _, v := range n.Pairs {
			k, _ := c.compileExpr(v.Key, false)
			c.rAlloc++
			v, sv := c.compileExpr(v.Value, false)
			switch sv {
			case rKonst:
				c.emitISetK(objAddr, k, v, 1)
			case rLoc:
				c.emitISetK(objAddr, k, v, 0)
			case rGlob:
				c.rAlloc++
				c.emitLoadG(v, c.rAlloc)
				c.emitISetK(objAddr, k, c.rAlloc, 0)
				c.rAlloc--
			case rFree:
				c.rAlloc++
				c.emitLoadF(v, c.rAlloc)
				c.emitISetK(objAddr, k, c.rAlloc, 0)
				c.rAlloc--
			}
			c.rAlloc--
		}
		return objAddr, rLoc
	case *ast.Property:
		return c.kb.StringIndex(n.Value), rKonst
	case *ast.ForState:
		return c.kb.IntegerIndex(0), rKonst
	case *ast.IGet:
		i, s := c.compileExpr(n.Indexable, false)
		lreg := c.rAlloc
		switch s {
		case rLoc:
			c.rAlloc++
			j, t := c.compileExpr(n.Index, false)
			switch t {
			case rLoc:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 0)
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.rAlloc--
				}
			case rKonst:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 1)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 1)
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.rAlloc--
				}
			}
		case rGlob:
			j, t := c.compileExpr(n.Index, false)
			switch t {
			case rLoc:
				c.rAlloc++
				c.emitLoadG(i, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(c.rAlloc, lreg, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(c.rAlloc, lreg, lreg, 0)
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			case rKonst:
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
				}
			case rFree:
				c.emitLoadG(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			}
		case rFree:
			j, t := c.compileExpr(n.Index, false)
			switch t {
			case rLoc:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadF(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
				} else {
					c.emitIGet(lreg, lreg+1, c.rAlloc, 0)
				}
			case rKonst:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
				}
			case rFree:
				c.emitLoadF(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			}
		}
		return c.rAlloc, rLoc
	case *ast.Select:
		i, s := c.compileExpr(n.Selectable, false)
		lreg := c.rAlloc
		switch s {
		case rLoc:
			c.rAlloc++
			j, t := c.compileExpr(n.Selector, false)
			switch t {
			case rLoc:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 0)
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.rAlloc--
				}
			case rKonst:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 1)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 1)
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.rAlloc--
				}
			}
		case rGlob:
			j, t := c.compileExpr(n.Selector, false)
			switch t {
			case rLoc:
				c.rAlloc++
				c.emitLoadG(i, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(c.rAlloc, lreg, c.rDest, 0)
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(c.rAlloc, lreg, lreg, 0)
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			case rKonst:
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
				}
			case rFree:
				c.emitLoadG(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			}
		case rFree:
			j, t := c.compileExpr(n.Selector, false)
			switch t {
			case rLoc:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, c.rAlloc, 0)
				}
			case rGlob:
				c.emitLoadF(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
				} else {
					c.emitIGet(lreg, lreg+1, c.rAlloc, 0)
				}
			case rKonst:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
				}
			case rFree:
				c.emitLoadF(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
				}
			}
		}
		return c.rAlloc, rLoc
	case *ast.Slice:
		v, s := c.compileExpr(n.Value, false)
		c.exprToReg(v, s)
		switch n.Mode {
		case vcv:
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
			}
		case vce:
			c.rAlloc++
			v, s := c.compileExpr(n.Last, false)
			c.exprToReg(v, s)
			c.rAlloc--
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
			}
		case ecv:
			c.rAlloc++
			v, s := c.compileExpr(n.First, false)
			c.exprToReg(v, s)
			c.rAlloc--
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
			}
		case ece:
			c.rAlloc++
			f, sf := c.compileExpr(n.First, false)
			c.exprToReg(f, sf)
			c.rAlloc++
			l, sl := c.compileExpr(n.Last, false)
			c.exprToReg(l, sl)
			c.rAlloc -= 2
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
			}
		}
		return c.rAlloc, rLoc
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
		idx, s := c.compileExpr(n.Fun, false)
		c.exprToReg(idx, s)
		for _, v := range n.Args {
			c.rAlloc++
			i, s := c.compileExpr(v, false)
			c.exprToReg(i, s)
		}
		c.rAlloc = reg
		c.emitCall(reg, len(n.Args))
		return reg, rLoc
	case *ast.MethodCallExpr:
		reg := c.rAlloc
		c.rAlloc++
		i, s := c.compileExpr(n.Obj, false)
		c.exprToReg(i, s)
		c.rAlloc++
		j, t := c.compileExpr(n.Prop, false)
		c.exprToReg(j, t)
		c.emitIGet(reg+1, reg+2, reg, 0)
		for _, v := range n.Args {
			i, s := c.compileExpr(v, false)
			c.exprToReg(i, s)
			c.rAlloc++
		}
		c.rAlloc = reg
		c.emitCall(reg, len(n.Args)+1)
		return reg, rLoc
	default:
		return 0, rGlob
	}
}

func (c *Compiler) compileConditional(n *ast.If, shouldJumpOutside bool) {
	idx, scope := c.compileExpr(n.Condition, false)
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
		c.exprToReg(idx, scope)
		addr := len(c.currentFn.Code)
		c.emitCheck(0, c.rAlloc, 0)
		c.compileBlockAndCheckJump(n.Block, shouldJumpOutside)
		c.currentFn.Code[addr] |= uint64(len(c.currentFn.Code))
	}
}

func (c *Compiler) skipBlock(block ast.Node) {
	addr := len(c.currentFn.Code)
	c.emitJump(0)
	c.compileStmt(block)
	c.currentFn.Code[addr] |= uint64(len(c.currentFn.Code))
}

func (c *Compiler) compileBlockAndCheckJump(block ast.Node, shouldJumpOutside bool) {
	c.compileStmt(block)
	if shouldJumpOutside {
		c.jumps = append(c.jumps, len(c.currentFn.Code))
		c.emitJump(0)
	}
}

func (c *Compiler) cleanUpLoopScope(init int, isWhileLoop bool) {
	hasBreaks := len(c.breakJumps)
	lastElem := len(c.breakCount) - 1
	count := c.breakCount[lastElem]
	if hasBreaks > 0 {
		for i := 1; i <= count; i++ {
			c.currentFn.Code[c.breakJumps[hasBreaks-i]] |= uint64(len(c.currentFn.Code))
		}
		c.breakJumps = c.breakJumps[:hasBreaks-count]
	}
	c.breakCount = c.breakCount[:lastElem]
	hasContinues := len(c.continueJumps)
	lastElem = len(c.continueCount) - 1
	count = c.continueCount[lastElem]
	if hasContinues > 0 {
		for i := 1; i <= count; i++ {
			if isWhileLoop {
				c.currentFn.Code[c.continueJumps[hasContinues-i]] |= uint64(init)
			} else {
				c.currentFn.Code[c.continueJumps[hasContinues-i]] |= uint64(len(c.currentFn.Code) - 1)
			}
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

func (c *Compiler) exprToReg(v, s int) {
	switch s {
	case rLoc:
		if v != c.rAlloc {
			c.emitMove(v, c.rAlloc)
		}
	case rGlob:
		c.emitLoadG(v, c.rAlloc)
	case rKonst:
		c.emitLoadK(v, c.rAlloc)
	case rFree:
		c.emitLoadF(v, c.rAlloc)
	}
}

func (c *Compiler) compileBinaryExpr(n *ast.BinaryExpr, isRoot bool) (int, int) {
	lidx, lscope := c.compileExpr(n.Lhs, false)
	lreg := c.rAlloc
	switch lscope {
	case rKonst:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rKonst:
			if val, err := c.kb.Konstants[lidx].Binop(uint64(n.Op), c.kb.Konstants[ridx]); err == nil {
				return c.integrateKonst(val)
			} else {
				c.hadError = true
			}
		case rGlob:
			c.emitLoadG(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, lreg, lreg, n.Op)
			}
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, ridx, lreg, n.Op)
			}
		case rFree:
			c.emitLoadF(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, lreg, lreg, n.Op)
			}
		}
	case rGlob:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rGlob:
			if c.mutLoc && isRoot {
				c.emitBinopG(lidx, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopG(lidx, ridx, lreg, n.Op)
			}
		case rKonst:
			c.emitLoadG(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lreg, lreg, n.Op)
			}
		case rLoc:
			c.rAlloc++
			c.emitLoadG(lidx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(c.rAlloc, lreg, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(c.rAlloc, lreg, lreg, n.Op)
				c.rAlloc--
			}
		case rFree:
			c.rAlloc++
			c.emitLoadG(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	case rLoc:
		c.rAlloc++
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitBinop(lidx, ridx, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, ridx, lreg, n.Op)
				c.rAlloc--
			}
		case rGlob:
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		case rKonst:
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lidx, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lidx, lreg, n.Op)
				c.rAlloc--
			}
		case rFree:
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	case rFree:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rLoc:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, ridx, lreg, n.Op)
			}
		case rGlob:
			c.rAlloc++
			c.emitLoadF(lidx, lreg)
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		case rKonst:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lreg, lreg, n.Op)
			}
		case rFree:
			c.rAlloc++
			c.emitLoadF(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	}
	return lreg, rLoc
}

func (c *Compiler) compileBinaryEq(n *ast.BinaryExpr, isRoot bool) (int, int) {
	lidx, lscope := c.compileExpr(n.Lhs, false)
	lreg := c.rAlloc
	switch lscope {
	case rKonst:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rKonst:
			val := c.kb.Konstants[lidx].Equals(c.kb.Konstants[ridx])
			return c.integrateKonst(val)
		case rGlob:
			c.emitLoadG(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitEqQ(lidx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqQ(lidx, lreg, lreg, n.Op)
			}
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitEqQ(lidx, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqQ(lidx, ridx, lreg, n.Op)
			}
		case rFree:
			c.emitLoadF(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitEqQ(lidx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqQ(lidx, lreg, lreg, n.Op)
			}
		}
	case rGlob:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rGlob:
			if c.mutLoc && isRoot {
				c.emitEqG(lidx, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqG(lidx, ridx, lreg, n.Op)
			}
		case rKonst:
			c.emitLoadG(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitEqK(ridx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqK(ridx, lreg, lreg, n.Op)
			}
		case rLoc:
			c.rAlloc++
			c.emitLoadG(lidx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(c.rAlloc, lreg, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(c.rAlloc, lreg, lreg, n.Op)
				c.rAlloc--
			}
		case rFree:
			c.rAlloc++
			c.emitLoadG(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	case rLoc:
		c.rAlloc++
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitEq(lidx, ridx, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lidx, ridx, lreg, n.Op)
				c.rAlloc--
			}
		case rGlob:
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(lidx, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lidx, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		case rKonst:
			if c.mutLoc && isRoot {
				c.emitEqK(ridx, lidx, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEqK(ridx, lidx, lreg, n.Op)
				c.rAlloc--
			}
		case rFree:
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(lidx, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lidx, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	case rFree:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rLoc:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitEq(lreg, ridx, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEq(lreg, ridx, lreg, n.Op)
			}
		case rGlob:
			c.rAlloc++
			c.emitLoadF(lidx, lreg)
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		case rKonst:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitEqK(ridx, lreg, c.rDest, n.Op)
				return c.rDest, rLoc
			} else {
				c.emitEqK(ridx, lreg, lreg, n.Op)
			}
		case rFree:
			c.rAlloc++
			c.emitLoadF(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitEq(lreg, c.rAlloc, c.rDest, n.Op)
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitEq(lreg, c.rAlloc, lreg, n.Op)
				c.rAlloc--
			}
		}
	}
	return lreg, rLoc
}
