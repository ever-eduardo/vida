package stdlib

import (
	"strings"

	"github.com/ever-eduardo/vida"
)

func generateText() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["hasPrefix"] = hasPrefix()
	m.Value["hasSuffix"] = hasSuffix()
	m.Value["fromCodepoint"] = fromCodepoint()
	m.Value["trim"] = trim()
	m.Value["trimLeft"] = trimLeft()
	m.Value["trimRight"] = trimRight()
	m.Value["split"] = split()
	m.UpdateKeys()
	return m
}

func hasPrefix() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return vida.Bool(strings.HasPrefix(v.Value, p.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func hasSuffix() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return vida.Bool(strings.HasSuffix(v.Value, p.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func fromCodepoint() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		runes := make([]rune, 0)
		for _, a := range args {
			if v, ok := a.(vida.Integer); ok && v > 0 {
				runes = append(runes, int32(v))
			}
		}
		return &vida.String{Value: string(runes), Runes: runes}, nil
	}
}

func trim() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return &vida.String{Value: strings.Trim(v.Value, p.Value)}, nil
				}
				return &vida.String{Value: strings.Trim(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return &vida.String{Value: strings.Trim(v.Value, " ")}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func trimLeft() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return &vida.String{Value: strings.TrimLeft(v.Value, p.Value)}, nil
				}
				return &vida.String{Value: strings.TrimLeft(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return &vida.String{Value: strings.TrimLeft(v.Value, " ")}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func trimRight() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return &vida.String{Value: strings.TrimRight(v.Value, p.Value)}, nil
				}
				return &vida.String{Value: strings.TrimRight(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return &vida.String{Value: strings.TrimRight(v.Value, " ")}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func split() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return stringSliceToList(strings.Split(v.Value, p.Value)), nil
				}
				return stringSliceToList(strings.Split(v.Value, "")), nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return stringSliceToList(strings.Split(v.Value, "")), nil
			}
		}
		return vida.NilValue, nil
	}
}

func stringSliceToList(slice []string) vida.Value {
	l := len(slice)
	xs := make([]vida.Value, l)
	for i := 0; i < l; i++ {
		xs[i] = &vida.String{Value: slice[i]}
	}
	return &vida.List{Value: xs}
}
