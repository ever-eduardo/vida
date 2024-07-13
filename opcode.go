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
	forInit
	forLoop
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
	forInit:  "ForInit",
	forLoop:  "ForLoop",
}
