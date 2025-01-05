package stdlib

import (
	"strconv"

	"github.com/ever-eduardo/vida"
)

func generateStringTo() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["boolean"] = toBool()
	m.Value["float"] = toFloat()
	m.Value["int"] = toInteger()
	m.Value["nilValue"] = toNil()
	m.UpdateKeys()
	return m
}

func toBool() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.String); ok {
				if v.Value == "true" {
					return vida.Bool(true), nil
				}
				if v.Value == "false" {
					return vida.Bool(false), nil
				}
			}
		}
		return vida.NilValue, nil
	}
}

func toNil() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.String); ok {
				if v.Value == "nil" {
					return vida.NilValue, nil
				}
			}
		}
		return vida.NilValue, nil
	}
}

func toFloat() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.String); ok {
				r, e := strconv.ParseFloat(v.Value, 64)
				if e == nil {
					return vida.Float(r), nil
				}
			}
		}
		return vida.NilValue, nil
	}
}

func toInteger() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.String); ok {
				i, e := strconv.ParseInt(v.Value, 0, 64)
				if e == nil {
					return vida.Integer(i), nil
				}
			}
		}
		return vida.NilValue, nil
	}
}
