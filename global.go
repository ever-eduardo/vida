package vida

import (
	"fmt"
)

var NilValue = Nil{}

func loadCoreLib() map[string]Value {
	p := make(map[string]Value)
	p["print"] = GoFn(gfnPrint)
	p["len"] = GoFn(gfnLen)
	p["append"] = NilValue
	p["load"] = NilValue
	p["type"] = NilValue
	p["assert"] = NilValue
	p["format"] = NilValue
	return p
}

func gfnPrint(args ...Value) (Value, error) {
	var s []any
	for _, v := range args {
		s = append(s, v.String())
	}
	fmt.Println(s...)
	return NilValue, nil
}

func gfnLen(args ...Value) (Value, error) {
	if len(args) > 0 {
		switch v := args[0].(type) {
		case *List:
			return Integer(len(v.Value)), nil
		case *Document:
			return Integer(len(v.Value)), nil
		case String:
			return Integer(len(v.Value)), nil
		}
	}
	return NilValue, nil
}

var strToRunesMap = make(map[string][]rune)
