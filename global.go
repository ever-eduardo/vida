package vida

var globalNil = Nil{}

var prelude = loadPrelude()

func loadPrelude() map[string]Value {
	p := make(map[string]Value)
	p["print"] = globalNil
	p["len"] = globalNil
	p["type"] = globalNil
	return p
}
