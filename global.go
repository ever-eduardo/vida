package vida

var globalNil = Nil{}

var prelude = loadPrelude()

func loadPrelude() map[string]Value {
	p := make(map[string]Value)
	p["print"] = "~Print"
	p["len"] = "~Len"
	p["type"] = "~Type"
	return p
}
