package vida

const (
	end = iota
	setG
	setL
	move
	prefix
	binop
	list
	iGet
	slice
)

var opcodes = [...]string{
	end:    "End",
	setG:   "SetG",
	setL:   "SetL",
	move:   "Move",
	prefix: "Prefix",
	binop:  "Binop",
	list:   "List",
	iGet:   "IGet",
	slice:  "Slice",
}
