package vida

type VM struct {
	Frame  [1024]frame
	Stack  [256]Value
	Module *Module
}

type frame struct {
	Code []byte
	Ip   int
	Op   byte
	Fp   int
}
