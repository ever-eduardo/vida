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
	obj
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
	obj:      "Obj",
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
