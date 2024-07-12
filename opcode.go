package vida

const (
	end = iota
	setG
	setL
	move
	prefix
	binop
	list
	document
	iGet
	iSet
	slice
)

var opcodes = [...]string{
	end:      "End",
	setG:     "SetG",
	setL:     "SetL",
	move:     "Move",
	prefix:   "Prefix",
	binop:    "Binop",
	list:     "List",
	document: "Doc",
	iGet:     "IGet",
	iSet:     "ISet",
	slice:    "Slice",
}
