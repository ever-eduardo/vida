package vida

const (
	end = iota
	setks
	locks
	move
	// Old School
	setAtom
	loadAtom
	setGlobal
	loadGlobal
	setLocal
)

const atomNil = 0
const atomTrue = 1
const atomFalse = 2

var opcodes = [...]string{
	end:        "End",
	setks:      "SetKS",
	locks:      "LocKS",
	move:       "Move",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
}
