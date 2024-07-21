package vida

import (
	"fmt"

	"github.com/ever-eduardo/vida/verror"
)

var NilValue = Nil{}

func loadCoreLib() map[string]Value {
	p := make(map[string]Value)
	p["print"] = GoFn(gfnPrint)
	p["len"] = GoFn(gfnLen)
	p["append"] = NilValue
	p["load"] = NilValue
	p["type"] = GoFn(gfnType)
	p["assert"] = GoFn(gfnAssert)
	p["format"] = GoFn(gfnFormat)
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

var strToRunesMap = make(map[string][]rune)
