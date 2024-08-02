package vida

import "github.com/ever-eduardo/vida/token"

const (
	rKonst = iota
	rLoc
	rGlob
	rFree
	rNotDefined
)

const (
	vcv = 2
	vce = 3
	ecv = 6
	ece = 7
)

const (
	shift16 = 16
	shift24 = 24
	shift32 = 32
	shift48 = 48
	shift56 = 56
	clean16 = 0x000000000000FFFF
	clean24 = 0x0000000000FFFFFF
)

func (c *Compiler) appendHeader() {
	c.currentFn.Code = append(c.currentFn.Code, header)
}

func (c *Compiler) appendEnd() {
	c.currentFn.Code = append(c.currentFn.Code, end)
}

func (c *Compiler) emitStoreG(from, to, isKonst int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= uint64(isKonst) << shift32
	i |= storeG << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitLoadG(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= loadG << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitLoadF(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= loadF << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitLoadK(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= loadK << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitMove(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= move << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitStoreF(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= storeF << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitLoc(from, to int, scope int) {
	// c.currentFn.Code = append(c.currentFn.Code, setL)
	// c.currentFn.Code = append(c.currentFn.Code, scope)
	// c.currentFn.Code = append(c.currentFn.Code, (from), (from>>8))
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitPrefix(from, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= uint64(operator) << shift32
	i |= prefix << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinop(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << 48
	i |= binop << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinopG(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << 48
	i |= binopG << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinopK(kidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(kidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << 48
	i |= binopK << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitEq(from, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= uint64(operator) << shift32
	i |= equals << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitList(length, root int) {
	var i uint64 = uint64(root)
	i |= uint64(length) << shift16
	i |= list << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitObject(length, from, to int) {
	// c.currentFn.Code = append(c.currentFn.Code, obj)
	// c.currentFn.Code = append(c.currentFn.Code, length)
	// c.currentFn.Code = append(c.currentFn.Code, from)
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitIGet(fromIndexable, fromIndex, scopeIndexable, scopeIndex, to int) {
	// c.currentFn.Code = append(c.currentFn.Code, iGet)
	// c.currentFn.Code = append(c.currentFn.Code, scopeIndexable)
	// c.currentFn.Code = append(c.currentFn.Code, scopeIndex)
	// c.currentFn.Code = append(c.currentFn.Code, (fromIndexable), (fromIndexable>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (fromIndex), (fromIndex>>8))
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitISet(fromIndex, fromExpr, scopeIndex, scopeExpr, from, to int) {
	// c.currentFn.Code = append(c.currentFn.Code, iSet)
	// c.currentFn.Code = append(c.currentFn.Code, scopeIndex)
	// c.currentFn.Code = append(c.currentFn.Code, scopeExpr)
	// c.currentFn.Code = append(c.currentFn.Code, (fromIndex), (fromIndex>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (fromExpr), (fromExpr>>8))
	// c.currentFn.Code = append(c.currentFn.Code, from)
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitSlice(mode, fromV, fromL, fromR, scopeV, scopeL, scopeR, to int) {
	// c.currentFn.Code = append(c.currentFn.Code, slice)
	// c.currentFn.Code = append(c.currentFn.Code, mode)
	// c.currentFn.Code = append(c.currentFn.Code, scopeV)
	// c.currentFn.Code = append(c.currentFn.Code, scopeL)
	// c.currentFn.Code = append(c.currentFn.Code, scopeR)
	// c.currentFn.Code = append(c.currentFn.Code, (fromV), (fromV>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (fromL), (fromL>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (fromR), (fromR>>8))
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitForSet(initReg, evalLoopAddr int) {
	// c.currentFn.Code = append(c.currentFn.Code, forSet)
	// c.currentFn.Code = append(c.currentFn.Code, initReg)
	// c.currentFn.Code = append(c.currentFn.Code, (evalLoopAddr), (evalLoopAddr>>8))
}

func (c *Compiler) emitForLoop(initReg, jump int) {
	// c.currentFn.Code = append(c.currentFn.Code, forLoop)
	// c.currentFn.Code = append(c.currentFn.Code, initReg)
	// c.currentFn.Code = append(c.currentFn.Code, (jump), (jump>>8))
}

func (c *Compiler) emitIForSet(evalLoopAddr, idx, scope, reg int) {
	// c.currentFn.Code = append(c.currentFn.Code, iForSet)
	// c.currentFn.Code = append(c.currentFn.Code, scope)
	// c.currentFn.Code = append(c.currentFn.Code, reg)
	// c.currentFn.Code = append(c.currentFn.Code, (idx), (idx>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (evalLoopAddr), (evalLoopAddr>>8))
}

func (c *Compiler) emitIForLoop(forLoopReg, jump int) {
	// c.currentFn.Code = append(c.currentFn.Code, iForLoop)
	// c.currentFn.Code = append(c.currentFn.Code, forLoopReg)
	// c.currentFn.Code = append(c.currentFn.Code, (jump), (jump>>8))
}

func (c *Compiler) emitJump(to int) {
	// c.currentFn.Code = append(c.currentFn.Code, jump)
	// c.currentFn.Code = append(c.currentFn.Code, (to), (to>>8))
}

func (c *Compiler) emitTestF(from, scope, jump int) {
	// c.currentFn.Code = append(c.currentFn.Code, checkF)
	// c.currentFn.Code = append(c.currentFn.Code, scope)
	// c.currentFn.Code = append(c.currentFn.Code, (from), (from>>8))
	// c.currentFn.Code = append(c.currentFn.Code, (jump), (jump>>8))
}

func (c *Compiler) emitFun(from, to int) {
	// c.currentFn.Code = append(c.currentFn.Code, fun)
	// c.currentFn.Code = append(c.currentFn.Code, (from), (fromtoken
	// c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitCall(from, argCount int) {
	// c.currentFn.Code = append(c.currentFn.Code, call)
	// c.currentFn.Code = append(c.currentFn.Code, from)
	// c.currentFn.Code = append(c.currentFn.Code, argCount)
}

func (c *Compiler) emitRet(from int, scope int) {
	// c.currentFn.Code = append(c.currentFn.Code, ret)
	// c.currentFn.Code = append(c.currentFn.Code, scope)
	// c.currentFn.Code = append(c.currentFn.Code, (from), (from>>8))
}

func (c *Compiler) refScope(id string) (int, int) {
	if to, isLocal, key := c.sb.isLocal(id); isLocal {
		if key.level != c.level {
			fn := c.fn[c.level]
			for i := 0; i < len(fn.Info); i++ {
				if fn.Info[i].Id == id {
					return i, rFree
				}
			}
			fn.Free++
			if key.level+1 == c.level {
				fn.Info = append(fn.Info, freeInfo{Index: int(to), IsLocal: true, Id: key.id})
			} else {
				for i := key.level + 1; i < c.level; i++ {
					if i == key.level+1 {
						c.fn[i].Free++
						c.fn[i].Info = append(c.fn[i].Info, freeInfo{Index: int(to), IsLocal: true, Id: key.id})
					} else {
						idx := len(c.fn[i-1].Info) - 1
						c.fn[i].Info = append(c.fn[i].Info, freeInfo{Index: idx, IsLocal: false, Id: key.id})
						c.fn[i].Free++
					}
				}
				fn.Info = append(fn.Info, freeInfo{Index: len(c.fn[c.level-1].Info) - 1, IsLocal: false, Id: key.id})
			}
			return len(fn.Info) - 1, rFree
		}
		return to, rLoc
	}
	if idx, isGlobal := c.sb.isGlobal(id); isGlobal {
		return idx, rGlob
	}
	c.hadError = true
	return 0, rNotDefined
}
