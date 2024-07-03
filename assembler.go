package vida

import "encoding/binary"

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

const (
	rKonst byte = iota
	rLocal
	rGlobal
	rPrelude
	rFreevar
)

const (
	opNot byte = iota
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

func (c *Compiler) refScope(id string) (int, byte) {
	if to, isLocal := c.sb.isLocal(id); isLocal {
		return int(to), rLocal
	}
	if isGlobal := c.sb.isGlobal(id); isGlobal {
		idx := c.kb.StringIndex(id)
		return idx, rGlobal
	}
	idx := c.kb.StringIndex(id)
	return idx, rPrelude
}
