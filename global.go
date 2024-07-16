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
	p["format"] = NilValue
	return p
}

func gfnPrint(args ...Value) (Value, error) {
	var s []any
	for i := range args {
		s = append(s, args[i].String(), ' ')
	}
	fmt.Println(s...)
	return NilValue, nil
}

var strToRunesMap = make(map[string][]rune)
