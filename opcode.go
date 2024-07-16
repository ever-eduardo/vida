package vida

const (
	end = iota
	setG
	setL
	move
	prefix
	binop
	equals
	list
	document
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
)

var opcodes = [...]string{
	end:      "End",
	setG:     "SetG",
	setL:     "SetL",
	move:     "Move",
	prefix:   "Prefix",
	binop:    "Binop",
	equals:   "Eq",
	list:     "List",
	document: "Doc",
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
}
