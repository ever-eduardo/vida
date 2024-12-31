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
	m.Value["fields"] = fields()
	m.Value["repeat"] = repeat()
	m.Value["replace"] = replace()
	m.Value["replaceAll"] = replaceAll()
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

func fields() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(vida.String); ok {
				return stringSliceToList(strings.Fields(v.Value)), nil
			}
		}
		return vida.NilValue, nil
	}
}

func repeat() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) >= 2 {
			if v, ok := args[0].(vida.String); ok {
				if times, ok := args[1].(vida.Integer); ok && times >= 0 {
					return &vida.String{Value: strings.Repeat(v.Value, int(times))}, nil
				}
				return vida.NilValue, nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func replace() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 3 {
			if s, ok := args[0].(vida.String); ok {
				if old, ok := args[1].(vida.String); ok {
					if nnew, ok := args[2].(vida.String); ok {
						if k, ok := args[3].(vida.Integer); ok {
							return &vida.String{Value: strings.Replace(s.Value, old.Value, nnew.Value, int(k))}, nil
						}
					}
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func replaceAll() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 2 {
			if s, ok := args[0].(vida.String); ok {
				if old, ok := args[1].(vida.String); ok {
					if nnew, ok := args[2].(vida.String); ok {
						return &vida.String{Value: strings.ReplaceAll(s.Value, old.Value, nnew.Value)}, nil
					}
				}
			}
			return vida.NilValue, nil
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
