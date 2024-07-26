package vida

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

const (
	rKonst byte = iota
	rLoc
	rGlob
	rFree
)

const (
	vcv = 2
	vce = 3
	ecv = 6
	ece = 7
)

func (c *Compiler) appendHeader() {
	c.currentFn.Code = append(c.currentFn.Code, byte(v))
	c.currentFn.Code = append(c.currentFn.Code, byte(i))
	c.currentFn.Code = append(c.currentFn.Code, byte(d))
	c.currentFn.Code = append(c.currentFn.Code, byte(a))
	c.currentFn.Code = append(c.currentFn.Code, byte(major))
	c.currentFn.Code = append(c.currentFn.Code, byte(minor))
	c.currentFn.Code = append(c.currentFn.Code, byte(patch))
	c.currentFn.Code = append(c.currentFn.Code, byte(inception))
}

func (c *Compiler) appendEnd() {
	c.currentFn.Code = append(c.currentFn.Code, end)
}

func (c *Compiler) emitSetG(from, to int, scope byte) {
	c.currentFn.Code = append(c.currentFn.Code, setG)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(to), byte(to>>8))
}

func (c *Compiler) emitLoc(from int, to byte, scope byte) {
	c.currentFn.Code = append(c.currentFn.Code, setL)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitSetF(from int, to byte, scope byte) {
	c.currentFn.Code = append(c.currentFn.Code, setF)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(uint16(to)), byte(uint16(to)>>8))
}

func (c *Compiler) emitMove(from byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, move)
	c.currentFn.Code = append(c.currentFn.Code, from)
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitPrefix(from int, to byte, scope byte, operator byte) {
	c.currentFn.Code = append(c.currentFn.Code, prefix)
	c.currentFn.Code = append(c.currentFn.Code, operator)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitBinary(fromLHS int, fromRHS int, scopeLHS byte, scopeRHS byte, to byte, operator byte) {
	c.currentFn.Code = append(c.currentFn.Code, binop)
	c.currentFn.Code = append(c.currentFn.Code, operator)
	c.currentFn.Code = append(c.currentFn.Code, scopeLHS)
	c.currentFn.Code = append(c.currentFn.Code, scopeRHS)
	c.currentFn.Code = append(c.currentFn.Code, byte(fromLHS), byte(fromLHS>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromRHS), byte(fromRHS>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitEq(fromLHS int, fromRHS int, scopeLHS byte, scopeRHS byte, to byte, operator byte) {
	c.currentFn.Code = append(c.currentFn.Code, equals)
	c.currentFn.Code = append(c.currentFn.Code, operator)
	c.currentFn.Code = append(c.currentFn.Code, scopeLHS)
	c.currentFn.Code = append(c.currentFn.Code, scopeRHS)
	c.currentFn.Code = append(c.currentFn.Code, byte(fromLHS), byte(fromLHS>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromRHS), byte(fromRHS>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitList(length byte, from byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, list)
	c.currentFn.Code = append(c.currentFn.Code, length)
	c.currentFn.Code = append(c.currentFn.Code, from)
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitObject(length byte, from byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, obj)
	c.currentFn.Code = append(c.currentFn.Code, length)
	c.currentFn.Code = append(c.currentFn.Code, from)
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitIGet(fromIndexable int, fromIndex int, scopeIndexable byte, scopeIndex byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, iGet)
	c.currentFn.Code = append(c.currentFn.Code, scopeIndexable)
	c.currentFn.Code = append(c.currentFn.Code, scopeIndex)
	c.currentFn.Code = append(c.currentFn.Code, byte(fromIndexable), byte(fromIndexable>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromIndex), byte(fromIndex>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitISet(fromIndex int, fromExpr int, scopeIndex byte, scopeExpr byte, from byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, iSet)
	c.currentFn.Code = append(c.currentFn.Code, scopeIndex)
	c.currentFn.Code = append(c.currentFn.Code, scopeExpr)
	c.currentFn.Code = append(c.currentFn.Code, byte(fromIndex), byte(fromIndex>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromExpr), byte(fromExpr>>8))
	c.currentFn.Code = append(c.currentFn.Code, from)
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitSlice(mode byte, fromV int, fromL int, fromR int, scopeV byte, scopeL byte, scopeR byte, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, slice)
	c.currentFn.Code = append(c.currentFn.Code, mode)
	c.currentFn.Code = append(c.currentFn.Code, scopeV)
	c.currentFn.Code = append(c.currentFn.Code, scopeL)
	c.currentFn.Code = append(c.currentFn.Code, scopeR)
	c.currentFn.Code = append(c.currentFn.Code, byte(fromV), byte(fromV>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromL), byte(fromL>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(fromR), byte(fromR>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitForSet(initReg byte, evalLoopAddr int) {
	c.currentFn.Code = append(c.currentFn.Code, forSet)
	c.currentFn.Code = append(c.currentFn.Code, initReg)
	c.currentFn.Code = append(c.currentFn.Code, byte(evalLoopAddr), byte(evalLoopAddr>>8))
}

func (c *Compiler) emitForLoop(initReg byte, jump int) {
	c.currentFn.Code = append(c.currentFn.Code, forLoop)
	c.currentFn.Code = append(c.currentFn.Code, initReg)
	c.currentFn.Code = append(c.currentFn.Code, byte(jump), byte(jump>>8))
}

func (c *Compiler) emitIForSet(evalLoopAddr, idx int, scope byte, reg byte) {
	c.currentFn.Code = append(c.currentFn.Code, iForSet)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, reg)
	c.currentFn.Code = append(c.currentFn.Code, byte(idx), byte(idx>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(evalLoopAddr), byte(evalLoopAddr>>8))
}

func (c *Compiler) emitIForLoop(forLoopReg byte, jump int) {
	c.currentFn.Code = append(c.currentFn.Code, iForLoop)
	c.currentFn.Code = append(c.currentFn.Code, forLoopReg)
	c.currentFn.Code = append(c.currentFn.Code, byte(jump), byte(jump>>8))
}

func (c *Compiler) emitJump(to int) {
	c.currentFn.Code = append(c.currentFn.Code, jump)
	c.currentFn.Code = append(c.currentFn.Code, byte(to), byte(to>>8))
}

func (c *Compiler) emitTestF(from int, scope byte, jump int) {
	c.currentFn.Code = append(c.currentFn.Code, testF)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, byte(jump), byte(jump>>8))
}

func (c *Compiler) emitFun(from int, to byte) {
	c.currentFn.Code = append(c.currentFn.Code, fun)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
	c.currentFn.Code = append(c.currentFn.Code, to)
}

func (c *Compiler) emitCall(from, argCount byte) {
	c.currentFn.Code = append(c.currentFn.Code, call)
	c.currentFn.Code = append(c.currentFn.Code, from)
	c.currentFn.Code = append(c.currentFn.Code, argCount)
}

func (c *Compiler) emitRet(from int, scope byte) {
	c.currentFn.Code = append(c.currentFn.Code, ret)
	c.currentFn.Code = append(c.currentFn.Code, scope)
	c.currentFn.Code = append(c.currentFn.Code, byte(from), byte(from>>8))
}

func (c *Compiler) refScope(id string) (int, byte) {
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
		return int(to), rLoc
	}
	if isGlobal := c.sb.isGlobal(id); isGlobal {
		idx := c.kb.StringIndex(id)
		return idx, rGlob
	}
	idx := c.kb.StringIndex(id)
	return idx, rGlob
}
