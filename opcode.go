package vida

const (
	end = iota
	setG
	setL
	move
	prefix
)

var opcodes = [...]string{
	end:    "End",
	setG:   "SetG",
	setL:   "SetL",
	move:   "Move",
	prefix: "Prefix",
}
