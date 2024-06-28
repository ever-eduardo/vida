package vida

import "encoding/binary"

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

const refKns byte = 0
const refStr byte = 1
const refReg byte = 3
const refPre byte = 4
const refFVr byte = 4

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

func (c *Compiler) emitSetSK(from, to int, flag byte) {
	c.module.Code = append(c.module.Code, setSK)
	c.module.Code = append(c.module.Code, flag)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(to))
}

func (c *Compiler) emitLocSK(from int, to byte, flag byte) {
	c.module.Code = append(c.module.Code, locSK)
	c.module.Code = append(c.module.Code, flag)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(from))
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) emitMove(from byte, to byte) {
	c.module.Code = append(c.module.Code, move)
	c.module.Code = append(c.module.Code, from)
	c.module.Code = append(c.module.Code, to)
}

func (c *Compiler) referenceScope(id string) (int, byte) {
	if iLoc, isLocal := c.lb.IsLocal(id, c.level, c.scope); isLocal {
		return int(iLoc), refReg
	} else if iGlob, isGlobal := c.kb.IsGlobal(id); isGlobal {
		return iGlob, refStr
	}
	return 0, refPre
}
