package vida

const (
	end = iota
	setKS
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
	setKS:      "SETKS",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
	readTop:    "ReadTop",
}
