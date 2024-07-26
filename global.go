package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/verror"
)

var NilValue = Nil{}

func loadCoreLib() map[string]Value {
	p := make(map[string]Value)
	p["print"] = GFn(gfnPrint)
	p["len"] = GFn(gfnLen)
	p["append"] = GFn(gfnAppend)
	p["make"] = GFn(gfnMakeArray)
	p["load"] = NilValue
	p["type"] = GFn(gfnType)
	p["assert"] = GFn(gfnAssert)
	p["fmt"] = GFn(gfnFormat)
	p["str"] = GFn(gfnString)
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
		case *Object:
			return Integer(len(v.Value)), nil
		case String:
			if v.Runes == nil {
				v.Runes = []rune(v.Value)
			}
			return Integer(len(v.Runes)), nil
		}
	}
	return NilValue, nil
}

func gfnType(args ...Value) (Value, error) {
	if len(args) > 0 {
		return String{Value: args[0].Type()}, nil
	}
	return NilValue, nil
}

func gfnFormat(args ...Value) (Value, error) {
	if len(args) > 1 {
		switch v := args[0].(type) {
		case String:
			s, e := formatValue(v.Value, args[1:]...)
			return String{Value: s}, e
		}
	}
	return NilValue, nil
}

func gfnAssert(args ...Value) (Value, error) {
	if len(args) > 0 {
		if args[0].Boolean() {
			return NilValue, nil
		} else {
			return NilValue, verror.AssertErr
		}
	}
	return NilValue, nil
}

func gfnAppend(args ...Value) (Value, error) {
	if len(args) >= 2 {
		switch v := args[0].(type) {
		case *List:
			v.Value = append(v.Value, args[1:]...)
			return v, nil
		}
	}
	return NilValue, nil
}

func gfnMakeArray(args ...Value) (Value, error) {
	largs := len(args)
	if largs > 0 {
		switch v := args[0].(type) {
		case Integer:
			var init Value = NilValue
			if largs > 1 {
				init = args[1]
			}
			if v >= 0 && v <= 0x7FFF_FFFF {
				arr := make([]Value, v)
				for i := range v {
					arr[i] = init
				}
				return &List{Value: arr}, nil
			}
		}
	}
	return NilValue, nil
}

func gfnString(args ...Value) (Value, error) {
	if len(args) > 0 {
		return String{Value: args[0].String()}, nil
	}
	return NilValue, nil
}
