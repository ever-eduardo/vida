package vida

const (
	end = iota
	setG
	setL
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
	setG:       "SetG",
	setL:       "SetL",
	move:       "Move",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
}
