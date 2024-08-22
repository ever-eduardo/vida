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
	i |= uint64(operator) << shift48
	i |= binop << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinopG(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << shift48
	i |= binopG << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinopK(kidx, regAddr, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(kidx) << shift16
	i |= uint64(regAddr) << shift32
	i |= uint64(operator) << shift48
	i |= binopK << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitBinopQ(kidx, regAddr, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(kidx) << shift16
	i |= uint64(regAddr) << shift32
	i |= uint64(operator) << shift48
	i |= binopQ << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitEq(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << shift48
	i |= eq << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitEqG(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << shift48
	i |= eqG << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitEqK(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << shift48
	i |= eqK << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitEqQ(lidx, ridx, to int, operator token.Token) {
	var i uint64 = uint64(to)
	i |= uint64(lidx) << shift16
	i |= uint64(ridx) << shift32
	i |= uint64(operator) << shift48
	i |= eqQ << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitList(length, root, to int) {
	var i uint64 = uint64(to)
	i |= uint64(root) << shift16
	i |= uint64(length) << shift32
	i |= list << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitObject(to int) {
	var i uint64 = uint64(to)
	i |= object << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitIGet(indexable, index, to, isKonst int) {
	var i uint64 = uint64(to)
	i |= uint64(index) << shift16
	i |= uint64(indexable) << shift32
	i |= uint64(isKonst) << shift48
	i |= iGet << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitISet(indexable, index, expr, isKonst int) {
	var i uint64 = uint64(expr)
	i |= uint64(index) << shift16
	i |= uint64(indexable) << shift32
	i |= uint64(isKonst) << shift48
	i |= iSet << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitISetK(indexable, index, expr, isKonst int) {
	var i uint64 = uint64(expr)
	i |= uint64(index) << shift16
	i |= uint64(indexable) << shift32
	i |= uint64(isKonst) << shift48
	i |= iSetK << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitSlice(mode, sliceable, to int) {
	var i uint64 = uint64(to)
	i |= uint64(sliceable) << shift16
	i |= uint64(mode) << shift32
	i |= slice << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitForSet(iReg, loop int) {
	var i uint64 = uint64(iReg)
	i |= uint64(loop) << shift16
	i |= forSet << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitForLoop(iReg, loop int) {
	var i uint64 = uint64(iReg)
	i |= uint64(loop) << shift16
	i |= forLoop << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitIForSet(loop, iterable, ireg int) {
	var i uint64 = uint64(ireg)
	i |= uint64(iterable) << shift16
	i |= uint64(loop) << shift32
	i |= iForSet << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitIForLoop(iReg, loop int) {
	var i uint64 = uint64(iReg)
	i |= uint64(loop) << shift16
	i |= iForLoop << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitJump(to int) {
	var i uint64 = uint64(to)
	i |= jump << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitCheck(against, reg, jump int) {
	var i uint64 = uint64(jump)
	i |= uint64(reg) << shift16
	i |= uint64(against) << shift32
	i |= check << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitFun(from, to int) {
	var i uint64 = uint64(to)
	i |= uint64(from) << shift16
	i |= fun << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitCall(callable, argCount, ellipsis, firstArg int) {
	var i uint64 = uint64(callable)
	i |= uint64(argCount) << shift16
	i |= uint64(ellipsis) << shift32
	i |= uint64(firstArg) << shift48
	i |= call << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
}

func (c *Compiler) emitRet(retReg int) {
	var i uint64 = uint64(retReg)
	i |= ret << shift56
	c.currentFn.Code = append(c.currentFn.Code, i)
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
