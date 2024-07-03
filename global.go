package vida

var NilValue = Nil{}

func loadPrelude() map[string]Value {
	p := make(map[string]Value)
	p["print"] = NilValue
	p["len"] = NilValue
	p["type"] = NilValue
	p["assert"] = NilValue
	return p
}
