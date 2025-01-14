package vida

const (
	end = iota
	load
	store
	prefix
	binop
	binopG
	binopK
	binopQ
	list
	object
	iGet
	iSet
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
	check
	fun
	ret
	call
)

var opcodes = [...]string{
	end:      "End",
	load:     "Load",
	store:    "Store",
	prefix:   "Prefix",
	binop:    "Binop",
	binopG:   "BinopG",
	binopK:   "BinopK",
	binopQ:   "BinopQ",
	list:     "List",
	object:   "Object",
	iGet:     "IGet",
	iSet:     "ISet",
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
	check:    "Check",
	fun:      "Fun",
	ret:      "Ret",
	call:     "Call",
}
