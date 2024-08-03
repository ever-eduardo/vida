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
	binop
	binopG
	binopK
	binopQ
	list
	object
	iGet
	iSet
	iSetK
	// --
	equals
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
	end:    "End",
	storeG: "StoreG",
	loadG:  "LoadG",
	loadF:  "LoadF",
	loadK:  "LoadK",
	move:   "Move",
	storeF: "StoreF",
	prefix: "Prefix",
	binop:  "Binop",
	binopG: "BinopG",
	binopK: "BinopK",
	binopQ: "BinopQ",
	list:   "List",
	object: "Object",
	iGet:   "IGet",
	iSet:   "ISet",
	iSetK:  "ISetK",
	// --
	equals:   "Eq",
	slice:    "Slice",
	forSet:   "For",
	forLoop:  "Loop",
	iForSet:  "IFor",
	iForLoop: "ILoop",
	checkF:   "CheckF",
	jump:     "Jump",
	fun:      "Fun",
	ret:      "Ret",
	call:     "Call",
}
