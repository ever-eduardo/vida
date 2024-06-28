package vida

const (
	end = iota
	setSK
	locSK
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
	setSK:      "SetSK",
	locSK:      "LocSK",
	move:       "Move",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
}
