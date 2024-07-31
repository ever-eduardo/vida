package vida

// 6 | 26
// 6 | 10 | 16
type Instruction uint32

const (
	end = iota
	setG
	setL
	setF
	move
	getR
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
	checkF
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
	getR:     "GetR",
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
	checkF:   "TestF",
	jump:     "Jump",
	fun:      "Fun",
	ret:      "Ret",
	call:     "Call",
}

func (i Instruction) String() string {
	return opcodes[i]
}
