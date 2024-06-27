package vida

import (
	"encoding/binary"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

const v rune = 'v'
const i rune = 'i'
const d rune = 'd'
const a rune = 'a'

var compilerError verror.VidaError

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

func (c *Compiler) makeNewGlobal() {
	c.compilationInfo.Let = false
	if c.compilationInfo.IsAtom {
		c.compilationInfo.IsAtom = false
		return
	}
}

func (c *Compiler) makeLocal() {
	c.compilationInfo.IsLocalAssignment = false
}

func (c *Compiler) makeAtomic(tok token.Token) {
	if c.compilationInfo.Let {
		c.compilationInfo.IsAtom = true
		c.module.Code = append(c.module.Code, setAtom)
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.symbolTable.AddGlobal(c.compilationInfo.Identifier))
	} else if c.compilationInfo.Var {
		c.module.Code = append(c.module.Code, loadAtom)
		c.module.Code = append(c.module.Code, c.rA)
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

func (c *Compiler) makeLoadRef() {
	if c.compilationInfo.Let {
		c.compilationInfo.IsAtom = true
		c.module.Code = append(c.module.Code, setGlobal)
		c.module.Code = binary.NativeEndian.AppendUint16(c.module.Code, c.symbolTable.AddGlobal(c.compilationInfo.Identifier))
	} else if c.compilationInfo.IsLocalAssignment {
		c.compilationInfo.IsAtom = true
	} else {
		c.module.Code = append(c.module.Code, loadGlobal)
		c.module.Code = append(c.module.Code, c.rA)
		c.rA++
	}
}
