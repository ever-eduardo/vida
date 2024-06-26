package vida

import (
	"encoding/binary"

	"github.com/ever-eduardo/vida/token"
)

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

func (c *Compiler) makeHeader() {
	c.module.Code = append(c.module.Code, byte(v))
	c.module.Code = append(c.module.Code, byte(i))
	c.module.Code = append(c.module.Code, byte(d))
	c.module.Code = append(c.module.Code, byte(a))
	c.module.Code = append(c.module.Code, byte(major))
	c.module.Code = append(c.module.Code, byte(minor))
	c.module.Code = append(c.module.Code, byte(patch))
	c.module.Code = append(c.module.Code, byte(inception))
}

func (c *Compiler) makeIdentifierPath() {
	c.compilationInfo.IsGlobalAssignment = false
	if c.compilationInfo.IsAtomicAssignment {
		c.compilationInfo.IsAtomicAssignment = false
		return
	}
}

func (c *Compiler) makeLocal(id string) {

}

func (c *Compiler) makeAtomic(tok token.Token) {
	if c.compilationInfo.IsGlobalAssignment {
		c.compilationInfo.IsAtomicAssignment = true
		c.module.Code = append(c.module.Code, setAtom)
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.konstantIndex(c.compilationInfo.Identifier))
	} else {
		c.module.Code = append(c.module.Code, loadAtom)
		c.module.Code = append(c.module.Code, c.rA)
		c.rA++
	}
	switch tok {
	case token.TRUE:
		c.module.Code = append(c.module.Code, atomTrue)
	case token.FALSE:
		c.module.Code = append(c.module.Code, atomFalse)
	default:
		c.module.Code = append(c.module.Code, atomNil)
	}
}

func (c *Compiler) makeStopRun() {
	c.module.Code = append(c.module.Code, stopRun)
}

func (c *Compiler) makeLoadGlobal() {
	if c.compilationInfo.IsGlobalAssignment {
		c.compilationInfo.IsAtomicAssignment = true
		c.module.Code = append(c.module.Code, loadGlobal)
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.konstantIndex(c.current.lit))
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.konstantIndex(c.compilationInfo.Identifier))
	} else {
		c.module.Code = append(c.module.Code, loadGlobal)
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.konstantIndex(c.current.lit))
		c.module.Code = append(c.module.Code, c.rA)
		c.rA++
	}
}

func (c *Compiler) konstantIndex(id string) uint16 {
	if idx, isPresent := c.identifiersMap[id]; isPresent {
		return idx
	} else {
		idx = c.kIndex
		c.identifiersMap[id] = idx
		c.module.Konstants = append(c.module.Konstants, id)
		c.kIndex++
		return idx
	}
}
