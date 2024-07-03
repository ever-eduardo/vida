package vida

const (
	end = iota
	setG
	setL
	move
	prefix
	binop
)

var opcodes = [...]string{
	end:    "End",
	setG:   "SetG",
	setL:   "SetL",
	move:   "Move",
	prefix: "Prefix",
	binop:  "Binary",
}
