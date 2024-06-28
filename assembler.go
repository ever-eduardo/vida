package vida

import "encoding/binary"

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

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

func (c *Compiler) emitSetFromK(dest, src int, flag byte) {
	c.module.Code = append(c.module.Code, setKS)
	c.module.Code = append(c.module.Code, flag)
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(dest))
	c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, uint16(src))
}
