package vida

import (
	"encoding/binary"
)

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

const (
	rKonst byte = iota
	rLoc
	rGlob
	rCore
	rFree
)

const (
	vcv = 2
	vce = 3
	ecv = 6
	ece = 7
)

func (c *Compiler) appendHeader() {
	c.module.Code = append(c.module.Code, byte(v))
	c.module.Code = append(c.module.Code, byte(i))
	c.module.Code = append(c.module.Code, byte(d))
	c.module.Code = append(c.module.Code, byte(a))
	c.module.Code = append(c.module.Code, byte(major))
	c.module.Code = append(c.module.Code, byte(minor))
	c.module.Code = append(c.module.Code, byte(patch))
	c.module.Code = append(c.module.Code, byte(inception))
}

func (c *Compiler) appendEnd() {
	c.module.Code = append(c.module.Code, end)
}

func (c *Compiler) emitSetG(from, to int, scope byte) {
	c.module.Code = append(c.module.Code, setG)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(to))
}

func (c *Compiler) emitLoc(from int, to byte, scope byte) {
	c.module.Code = append(c.module.Code, setL)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitSetF(from int, to byte, scope byte) {
	c.module.Code = append(c.module.Code, setF)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(to))
}

func (c *Compiler) emitMove(from byte, to byte) {
	c.module.Code = append(c.module.Code, move)
	c.module.Code = append(c.module.Code, from)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitPrefix(from int, to byte, scope byte, operator byte) {
	c.module.Code = append(c.module.Code, prefix)
	c.module.Code = append(c.module.Code, operator)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitBinary(fromLHS int, fromRHS int, scopeLHS byte, scopeRHS byte, to byte, operator byte) {
	c.module.Code = append(c.module.Code, binop)
	c.module.Code = append(c.module.Code, operator)
	c.module.Code = append(c.module.Code, scopeLHS)
	c.module.Code = append(c.module.Code, scopeRHS)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromLHS))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromRHS))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitEq(fromLHS int, fromRHS int, scopeLHS byte, scopeRHS byte, to byte, operator byte) {
	c.module.Code = append(c.module.Code, equals)
	c.module.Code = append(c.module.Code, operator)
	c.module.Code = append(c.module.Code, scopeLHS)
	c.module.Code = append(c.module.Code, scopeRHS)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromLHS))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromRHS))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitList(length byte, from byte, to byte) {
	c.module.Code = append(c.module.Code, list)
	c.module.Code = append(c.module.Code, length)
	c.module.Code = append(c.module.Code, from)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitDocument(length byte, from byte, to byte) {
	c.module.Code = append(c.module.Code, doc)
	c.module.Code = append(c.module.Code, length)
	c.module.Code = append(c.module.Code, from)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitIGet(fromIndexable int, fromIndex int, scopeIndexable byte, scopeIndex byte, to byte) {
	c.module.Code = append(c.module.Code, iGet)
	c.module.Code = append(c.module.Code, scopeIndexable)
	c.module.Code = append(c.module.Code, scopeIndex)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromIndexable))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromIndex))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitISet(fromIndex int, fromExpr int, scopeIndex byte, scopeExpr byte, from byte, to byte) {
	c.module.Code = append(c.module.Code, iSet)
	c.module.Code = append(c.module.Code, scopeIndex)
	c.module.Code = append(c.module.Code, scopeExpr)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromIndex))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromExpr))
	c.module.Code = append(c.module.Code, from)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitSlice(mode byte, fromV int, fromL int, fromR int, scopeV byte, scopeL byte, scopeR byte, to byte) {
	c.module.Code = append(c.module.Code, slice)
	c.module.Code = append(c.module.Code, mode)
	c.module.Code = append(c.module.Code, scopeV)
	c.module.Code = append(c.module.Code, scopeL)
	c.module.Code = append(c.module.Code, scopeR)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromV))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromL))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(fromR))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitForSet(initReg byte, evalLoopAddr int) {
	c.module.Code = append(c.module.Code, forSet)
	c.module.Code = append(c.module.Code, initReg)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(evalLoopAddr))
}

func (c *Compiler) emitForLoop(initReg byte, jump int) {
	c.module.Code = append(c.module.Code, forLoop)
	c.module.Code = append(c.module.Code, initReg)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(jump))
}

func (c *Compiler) emitIForSet(evalLoopAddr, idx int, scope byte, reg byte) {
	c.module.Code = append(c.module.Code, iForSet)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = append(c.module.Code, reg)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(idx))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(evalLoopAddr))
}

func (c *Compiler) emitIForLoop(forLoopReg byte, jump int) {
	c.module.Code = append(c.module.Code, iForLoop)
	c.module.Code = append(c.module.Code, forLoopReg)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(jump))
}

func (c *Compiler) emitJump(to int) {
	c.module.Code = append(c.module.Code, jump)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(to))
}

func (c *Compiler) emitTestF(from int, scope byte, jump int) {
	c.module.Code = append(c.module.Code, testF)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(jump))
}

func (c *Compiler) emitFun(from int, to byte, jump int) {
	c.module.Code = append(c.module.Code, fun)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(jump))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitCall(fn byte, argCount, to byte) {
	c.module.Code = append(c.module.Code, call)
	c.module.Code = append(c.module.Code, fn)
	c.module.Code = append(c.module.Code, argCount)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitRet(from int, scope byte) {
	c.module.Code = append(c.module.Code, ret)
	c.module.Code = append(c.module.Code, scope)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
}

func (c *Compiler) refScope(id string) (int, byte) {
	if to, isLocal, key := c.sb.isLocal(id); isLocal {
		if key.level != c.level {
			fn := c.fn[c.level-1]
			for i := 0; i < len(fn.Info); i++ {
				if fn.Info[i].Id == id {
					return i, rFree
				}
			}
			fn.Free++
			if key.level+1 == c.level {
				fn.Info = append(fn.Info, freeInfo{Index: int(to), IsLocal: true, Id: key.id})
			} else {
				for i := key.level; i < c.level-1; i++ {
					if i == key.level {
						c.fn[i].Free++
						c.fn[i].Info = append(c.fn[i].Info, freeInfo{Index: int(to), IsLocal: true, Id: key.id})
					} else {
						idx := len(c.fn[i-1].Info) - 1
						c.fn[i].Info = append(c.fn[i].Info, freeInfo{Index: idx, IsLocal: false, Id: key.id})
						c.fn[i].Free++
					}
				}
				fn.Info = append(fn.Info, freeInfo{Index: len(c.fn[c.level-2].Info) - 1, IsLocal: false, Id: key.id})
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
	return idx, rCore
}
