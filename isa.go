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
	eq
	eqG
	eqK
	eqQ
	slice
	forSet
	forLoop
	jump
	iForSet
	iForLoop
	// --
	checkF
	fun
	ret
	call
)

var opcodes = [...]string{
	end:      "End",
	storeG:   "StoreG",
	loadG:    "LoadG",
	loadF:    "LoadF",
	loadK:    "LoadK",
	move:     "Move",
	storeF:   "StoreF",
	prefix:   "Prefix",
	binop:    "Binop",
	binopG:   "BinopG",
	binopK:   "BinopK",
	binopQ:   "BinopQ",
	list:     "List",
	object:   "Object",
	iGet:     "IGet",
	iSet:     "ISet",
	iSetK:    "ISetK",
	eq:       "Eq",
	eqG:      "EqG",
	eqK:      "EqK",
	eqQ:      "EqQ",
	slice:    "Slice",
	forSet:   "For",
	forLoop:  "Loop",
	jump:     "Jump",
	iForSet:  "IFor",
	iForLoop: "ILoop",
	// --
	checkF: "CheckF",
	fun:    "Fun",
	ret:    "Ret",
	call:   "Call",
}
