package vida

import (
	"github.com/ever-eduardo/vida/ast"
	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type compiler struct {
	jumps         []int
	breakJumps    []int
	breakCount    []int
	continueJumps []int
	continueCount []int
	errMsg        string
	fn            []*CoreFunction
	currentFn     *CoreFunction
	ast           *ast.Ast
	module        *Module
	kb            *konstBuilder
	sb            *symbolBuilder
	moduleMap     map[string]int
	depMap        map[string]struct{}
	linesMap      map[string]map[int]uint
	lineErr       uint
	scope         int
	level         int
	rAlloc        int
	rDest         int
	fromRefStmt   bool
	mutLoc        bool
	hadError      bool
	isSubcompiler bool
}

var dummy = struct{}{}

func newMainCompiler(ast *ast.Ast, moduleName string) *compiler {
	dm := make(map[string]struct{})
	dm[moduleName] = dummy
	lm := make(map[string]map[int]uint)
	lm[moduleName] = make(map[int]uint)
	c := &compiler{
		ast:       ast,
		module:    newMainModule(moduleName),
		kb:        newKonstBuilder(),
		sb:        newSymbolBuilder(0),
		moduleMap: make(map[string]int),
		depMap:    dm,
		linesMap:  lm,
	}
	c.fn = append(c.fn, c.module.MainFunction.CoreFn)
	c.currentFn = c.module.MainFunction.CoreFn
	return c
}

func newSubCompiler(ast *ast.Ast, moduleName string, kb *konstBuilder, store *[]Value, moduleMap map[string]int, depMap map[string]struct{}, lm map[string]map[int]uint, initialIndex int) *compiler {
	lm[moduleName] = make(map[int]uint)
	c := &compiler{
		ast:           ast,
		module:        newSubModule(moduleName, store),
		kb:            kb,
		sb:            newSymbolBuilder(initialIndex),
		isSubcompiler: true,
		moduleMap:     moduleMap,
		depMap:        depMap,
		linesMap:      lm,
	}
	c.fn = append(c.fn, c.module.MainFunction.CoreFn)
	c.currentFn = c.module.MainFunction.CoreFn
	return c
}

func (c *compiler) compileModule() (*Module, error) {
	c.appendHeader()
	var i int
	for i = range len(c.ast.Statement) {
		c.compileStmt(c.ast.Statement[i])
		if c.hadError {
			return nil, verror.New(c.module.MainFunction.CoreFn.ModuleName, c.errMsg, verror.CompilationErrType, c.lineErr)
		}
	}
	c.module.Konstants = c.kb.Konstants
	c.appendEnd()
	return c.module, nil
}

func (c *compiler) compileSubModule() (*Module, error) {
	for i := range len(c.ast.Statement) {
		c.compileStmt(c.ast.Statement[i])
		if c.hadError {
			return nil, verror.New(c.module.MainFunction.CoreFn.ModuleName, c.errMsg, verror.CompilationErrType, c.lineErr)
		}
	}
	return c.module, nil
}

func (c *compiler) compileStmt(node ast.Node) {
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
		case rNotDefined:
			c.generateReferenceError(n.Indentifier, n.Line)
		}
	case *ast.Let:
		to, isPresent := c.sb.addGlobal(n.Indentifier)
		if !isPresent {
			*c.module.Store = append(*c.module.Store, NilValue)
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
		var from, scope int
		if n.IsRecursive {
			c.sb.addLocal(n.Identifier, c.level, c.scope, to)
			c.emitLoadK(c.kb.NilIndex(), to)
			from, scope = c.compileExpr(n.Expr, true)
		} else {
			from, scope = c.compileExpr(n.Expr, true)
			c.sb.addLocal(n.Identifier, c.level, c.scope, to)
		}
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
			switch v := (*c.kb.Konstants)[idx].(type) {
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
			c.exprToReg(idx, scope)
			addr := len(c.currentFn.Code)
			c.emitCheck(0, c.rAlloc, 0)
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
		case rNotDefined:
			c.generateReferenceError(n.Value, n.Line)
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
		if c.level != 0 || c.isSubcompiler {
			i, s := c.compileExpr(n.Expr, true)
			c.exprToReg(i, s)
			c.emitRet(c.rAlloc)
		}
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
		c.emitCall(callable, len(n.Args), n.Ellipsis, 1)
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		for _, v := range n.Args {
			i, s := c.compileExpr(v, true)
			c.exprToReg(i, s)
			c.rAlloc++
		}
		c.rAlloc = obj
		c.fromRefStmt = false
		c.emitCall(obj, len(n.Args)+1, n.Ellipsis, 2)
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
	case *ast.Export:
		if c.isSubcompiler {
			i, s := c.compileExpr(n.Expr, true)
			c.exprToReg(i, s)
			c.emitRet(c.rAlloc)
		}
	}
}

func (c *compiler) compileExpr(node ast.Node, isRoot bool) (int, int) {
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
			if val, err := (*c.kb.Konstants)[from].Prefix(uint64(n.Op)); err == nil {
				return c.integrateKonst(val)
			} else {
				c.hadError = true
				c.errMsg = "cannot perform prefix operation"
				c.lineErr = n.Line
			}
		}
		switch scope {
		case rGlob:
			c.emitLoadG(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitPrefix(from, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitPrefix(from, c.rAlloc, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case rKonst:
			c.emitLoadK(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		case rFree:
			c.emitLoadF(from, c.rAlloc)
			c.emitPrefix(c.rAlloc, c.rAlloc, n.Op)
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		}
		return c.rAlloc, rLoc
	case *ast.Boolean:
		return c.kb.BooleanIndex(n.Value), rKonst
	case *ast.Nil:
		return c.kb.NilIndex(), rKonst
	case *ast.Reference:
		i, s := c.refScope(n.Value)
		if s == rNotDefined {
			c.generateReferenceError(n.Value, n.Line)
		}
		return i, s
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
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		case rGlob:
			c.rAlloc++
			j, t := c.compileExpr(n.Index, false)
			switch t {
			case rLoc:
				c.emitMove(j, c.rAlloc)
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadG(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		case rFree:
			c.rAlloc++
			j, t := c.compileExpr(n.Index, false)
			switch t {
			case rLoc:
				c.emitMove(j, c.rAlloc)
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadF(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		}
		return lreg, rLoc
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
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				if c.mutLoc && isRoot {
					c.emitIGet(i, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(j, c.rAlloc)
				if c.mutLoc && isRoot {
					c.emitIGet(i, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(i, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		case rGlob:
			c.rAlloc++
			j, t := c.compileExpr(n.Selector, false)
			switch t {
			case rLoc:
				c.emitMove(j, c.rAlloc)
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadG(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				c.emitLoadG(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadG(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		case rFree:
			c.rAlloc++
			j, t := c.compileExpr(n.Selector, false)
			switch t {
			case rLoc:
				c.emitMove(j, c.rAlloc)
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, c.rAlloc, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, c.rAlloc, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rGlob:
				c.emitLoadF(i, lreg)
				c.emitLoadG(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rKonst:
				c.emitLoadF(i, lreg)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, j, c.rDest, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, j, lreg, 1)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			case rFree:
				c.emitLoadF(i, lreg)
				c.emitLoadF(j, lreg+1)
				if c.mutLoc && isRoot {
					c.emitIGet(lreg, lreg+1, c.rDest, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
					return c.rDest, rLoc
				} else {
					c.emitIGet(lreg, lreg+1, lreg, 0)
					c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
					c.rAlloc--
				}
			}
		}
		return lreg, rLoc
	case *ast.Slice:
		v, s := c.compileExpr(n.Value, false)
		c.exprToReg(v, s)
		switch n.Mode {
		case vcv:
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case vce:
			c.rAlloc++
			v, s := c.compileExpr(n.Last, false)
			c.exprToReg(v, s)
			c.rAlloc--
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case ecv:
			c.rAlloc++
			v, s := c.compileExpr(n.First, false)
			c.exprToReg(v, s)
			c.rAlloc--
			if c.mutLoc && isRoot {
				c.emitSlice(n.Mode, c.rAlloc, c.rDest)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitSlice(n.Mode, c.rAlloc, c.rAlloc)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		}
		return c.rAlloc, rLoc
	case *ast.Fun:
		fn := &CoreFunction{ModuleName: c.module.MainFunction.CoreFn.ModuleName}
		c.fn = append(c.fn, fn)
		c.emitFun(c.kb.FunctionIndex(fn), c.rAlloc)
		c.currentFn = fn
		reg := c.startFuncScope()
		for _, v := range n.Args {
			fn.Arity++
			c.sb.addLocal(v, c.level, c.scope, c.rAlloc)
			c.rAlloc++
		}
		if n.IsVar {
			fn.IsVar = true
			fn.Arity--
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
		c.emitCall(reg, len(n.Args), n.Ellipsis, 1)
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		for _, v := range n.Args {
			i, s := c.compileExpr(v, false)
			c.exprToReg(i, s)
			c.rAlloc++
		}
		c.rAlloc = reg
		c.emitCall(reg, len(n.Args)+1, n.Ellipsis, 2)
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		return reg, rLoc
	case *ast.Import:
		if _, isCycle := c.depMap[n.Path]; isCycle {
			c.hadError = true
			c.errMsg = "import cycle detected"
			c.lineErr = n.Line
			return 0, rGlob
		} else {
			c.depMap[n.Path] = dummy
		}
		if v, isPresent := c.moduleMap[n.Path]; isPresent {
			delete(c.depMap, n.Path)
			c.emitFun(v, c.rAlloc)
			c.emitCall(c.rAlloc, 0, 0, 1)
			c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			return c.rAlloc, rLoc
		}
		src, err := readModule(n.Path)
		if err != nil {
			c.hadError = true
			c.errMsg = err.Error()
			c.lineErr = n.Line
			return 0, rGlob
		}
		p := newParser(src, n.Path)
		moduleAST, err := p.parse()
		if err != nil {
			c.hadError = true
			c.errMsg = err.Error()
			c.lineErr = n.Line
			return 0, rGlob
		}
		subCompiler := newSubCompiler(moduleAST, n.Path, c.kb, c.module.Store, c.moduleMap, c.depMap, c.linesMap, len(*c.module.Store))
		m, err := subCompiler.compileSubModule()
		c.sb.index = len(*c.module.Store)
		if err != nil {
			c.hadError = true
			c.errMsg = err.Error()
			c.lineErr = n.Line
			return 0, rGlob
		}
		fnIndex := c.kb.FunctionIndex(m.MainFunction.CoreFn)
		c.moduleMap[n.Path] = fnIndex
		delete(c.depMap, n.Path)
		c.emitFun(fnIndex, c.rAlloc)
		c.emitCall(c.rAlloc, 0, 0, 1)
		c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
		return c.rAlloc, rLoc
	case *ast.Enum:
		e := make(Enum)
		if n.HasInitVal {
			for _, v := range n.Variants {
				if v == "_" {
					n.Init++
					continue
				}
				e[v] = Integer(n.Init)
				n.Init++
			}
			return c.kb.EnumIndex(e), rKonst
		}
		for i, v := range n.Variants {
			if v == "_" {
				continue
			}
			e[v] = Integer(i)
		}
		return c.kb.EnumIndex(e), rKonst
	default:
		return 0, rGlob
	}
}

func (c *compiler) compileConditional(n *ast.If, shouldJumpOutside bool) {
	idx, scope := c.compileExpr(n.Condition, false)
	if scope == rKonst {
		switch v := (*c.kb.Konstants)[idx].(type) {
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

func (c *compiler) skipBlock(block ast.Node) {
	addr := len(c.currentFn.Code)
	c.emitJump(0)
	c.compileStmt(block)
	c.currentFn.Code[addr] |= uint64(len(c.currentFn.Code))
}

func (c *compiler) compileBlockAndCheckJump(block ast.Node, shouldJumpOutside bool) {
	c.compileStmt(block)
	if shouldJumpOutside {
		c.jumps = append(c.jumps, len(c.currentFn.Code))
		c.emitJump(0)
	}
}

func (c *compiler) cleanUpLoopScope(init int, isWhileLoop bool) {
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

func (c *compiler) startLoopScope() {
	c.breakCount = append(c.breakCount, 0)
	c.continueCount = append(c.continueCount, 0)
}

func (c *compiler) startFuncScope() int {
	r := c.rAlloc
	c.rAlloc = 0
	c.level++
	return r
}

func (c *compiler) leaveFuncScope() {
	c.sb.clearLocals(c.level, c.scope)
	c.fn = c.fn[:c.level]
	c.level--
	c.currentFn = c.fn[c.level]
}

func (c *compiler) integrateKonst(val Value) (int, int) {
	switch e := val.(type) {
	case Integer:
		return c.kb.IntegerIndex(int64(e)), rKonst
	case Float:
		return c.kb.FloatIndex(float64(e)), rKonst
	case Bool:
		return c.kb.BooleanIndex(bool(e)), rKonst
	case *String:
		return c.kb.StringIndex(e.Value), rKonst
	default:
		return c.kb.NilIndex(), rKonst
	}
}

func (c *compiler) exprToReg(i, s int) {
	switch s {
	case rLoc:
		if i != c.rAlloc {
			c.emitMove(i, c.rAlloc)
		}
	case rGlob:
		c.emitLoadG(i, c.rAlloc)
	case rKonst:
		c.emitLoadK(i, c.rAlloc)
	case rFree:
		c.emitLoadF(i, c.rAlloc)
	}
}

func (c *compiler) compileBinaryExpr(n *ast.BinaryExpr, isRoot bool) (int, int) {
	lidx, lscope := c.compileExpr(n.Lhs, false)
	lreg := c.rAlloc
	switch lscope {
	case rKonst:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rKonst:
			if val, err := (*c.kb.Konstants)[lidx].Binop(uint64(n.Op), (*c.kb.Konstants)[ridx]); err == nil {
				return c.integrateKonst(val)
			} else {
				c.hadError = true
				c.errMsg = "cannot perform binary operation"
				c.lineErr = n.Line
			}
		case rGlob:
			c.emitLoadG(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, lreg, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, lreg, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case rLoc:
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, ridx, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, ridx, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case rFree:
			c.emitLoadF(ridx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopQ(lidx, lreg, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitBinopQ(lidx, lreg, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		}
	case rGlob:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rGlob:
			if c.mutLoc && isRoot {
				c.emitBinopG(lidx, ridx, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitBinopG(lidx, ridx, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case rKonst:
			c.emitLoadG(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lreg, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lreg, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
			}
		case rLoc:
			c.emitMove(ridx, lreg)
			c.rAlloc++
			c.emitLoadG(lidx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(c.rAlloc, lreg, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(c.rAlloc, lreg, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rFree:
			c.rAlloc++
			c.emitLoadG(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
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
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, ridx, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rGlob:
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rKonst:
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lidx, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lidx, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rFree:
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lidx, c.rAlloc, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lidx, c.rAlloc, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		}
	case rFree:
		c.rAlloc++
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rLoc:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, ridx, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, ridx, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rGlob:
			c.emitLoadF(lidx, lreg)
			c.emitLoadG(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rKonst:
			c.emitLoadF(lidx, lreg)
			if c.mutLoc && isRoot {
				c.emitBinopK(ridx, lreg, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinopK(ridx, lreg, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		case rFree:
			c.emitLoadF(lidx, lreg)
			c.emitLoadF(ridx, c.rAlloc)
			if c.mutLoc && isRoot {
				c.emitBinop(lreg, c.rAlloc, c.rDest, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
				return c.rDest, rLoc
			} else {
				c.emitBinop(lreg, c.rAlloc, lreg, n.Op)
				c.linesMap[c.currentFn.ModuleName][len(c.currentFn.Code)] = n.Line
				c.rAlloc--
			}
		}
	}
	return lreg, rLoc
}

func (c *compiler) compileBinaryEq(n *ast.BinaryExpr, isRoot bool) (int, int) {
	lidx, lscope := c.compileExpr(n.Lhs, false)
	lreg := c.rAlloc
	switch lscope {
	case rKonst:
		ridx, rscope := c.compileExpr(n.Rhs, false)
		switch rscope {
		case rKonst:
			val := (*c.kb.Konstants)[lidx].Equals((*c.kb.Konstants)[ridx])
			if n.Op == token.NEQ {
				val = !val
			}
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
			c.emitMove(ridx, lreg)
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
