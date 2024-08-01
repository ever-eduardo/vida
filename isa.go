package vida

const (
	// Refactoring
	end = iota
	storeG
	loadG
	loadF
	loadK
	move
	storeF
	prefix
	// --
	binop
	equals
	list
	obj
	iGet
	iSet
	slice
	forSet
	forLoop
	iForSet
	iForLoop
	checkF
	jump
	fun
	ret
	call
)

var opcodes = [...]string{
	end:      "End",
	storeG:   "StoreG",
	storeF:   "SetF",
	move:     "Move",
	loadG:    "LoadG",
	loadF:    "LoadF",
	loadK:    "LoadK",
	prefix:   "Prefix",
	binop:    "Binop",
	equals:   "Eq",
	list:     "List",
	obj:      "Obj",
	iGet:     "IGet",
	iSet:     "ISet",
	slice:    "Slice",
	forSet:   "For",
	forLoop:  "Loop",
	iForSet:  "IFor",
	iForLoop: "ILoop",
	checkF:   "TestF",
	jump:     "Jump",
	fun:      "Fun",
	ret:      "Ret",
	call:     "Call",
}
