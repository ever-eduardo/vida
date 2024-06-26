package vida

const (
	stopRun = iota
	setAtom
	loadAtom
	loadGlobal
)

const atomNil = 0
const atomTrue = 1
const atomFalse = 2

var opcodes = [...]string{
	stopRun:    "StopRun",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	loadGlobal: "LoadGlobal",
}
