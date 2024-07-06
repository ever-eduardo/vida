package vida

import (
	"fmt"
)

var NilValue = Nil{}

func loadPrelude() map[string]Value {
	p := make(map[string]Value)
	p["print"] = GoFn(gfnPrint)
	p["len"] = NilValue
	p["append"] = NilValue
	p["load"] = NilValue
	p["type"] = NilValue
	p["assert"] = NilValue
	return p
}

func gfnPrint(args ...Value) (Value, error) {
	for i := range args {
		fmt.Println(args[i].String())
	}
	return NilValue, nil
}
