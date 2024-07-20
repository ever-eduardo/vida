package vida

const (
	end = iota
	setG
	setL
	setF
	move
	prefix
	binop
	equals
	list
	doc
	iGet
	iSet
	slice
	forSet
	forLoop
	iForSet
	iForLoop
	testF
	test
	jump
	fun
	ret
	call
)

var opcodes = [...]string{
	end:      "End",
	setG:     "SetG",
	setL:     "SetL",
	setF:     "SetF",
	move:     "Move",
	prefix:   "Prefix",
	binop:    "Binop",
	equals:   "Eq",
	list:     "List",
	doc:      "Doc",
	iGet:     "IGet",
	iSet:     "ISet",
	slice:    "Slice",
	forSet:   "For",
	forLoop:  "Loop",
	iForSet:  "IFor",
	iForLoop: "ILoop",
	testF:    "TestF",
	test:     "Test",
	jump:     "Jump",
	fun:      "Fun",
	ret:      "Ret",
	call:     "Call",
}
