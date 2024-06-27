package vida

const (
	stopRun = iota
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
	stopRun:    "StopRun",
	setAtom:    "SetAtom",
	loadAtom:   "LoadAtom",
	setGlobal:  "SetGlobal",
	loadGlobal: "LoadGlobal",
	setLocal:   "SetLocal",
	readTop:    "ReadTop",
}
