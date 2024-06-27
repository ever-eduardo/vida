package vida

const (
	end = iota
	setK
	// Old School
	setAtom
	loadAtom
	setGlobal
	loadGlobal
	setLocal
	readTop
)

const atomNil = 0
const atomTrue = 1
const atomFalse = 2

var opcodes = [...]string{
	end:        "END",
	setK:       "SETK",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
	readTop:    "ReadTop",
}
