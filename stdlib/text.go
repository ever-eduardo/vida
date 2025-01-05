package stdlib

import (
	"strings"
	"unicode"
	"unicode/utf8"

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
	m.Value["center"] = center()
	m.Value["contains"] = contains()
	m.Value["containsAny"] = containsAny()
	m.Value["index"] = index()
	m.Value["join"] = join()
	m.Value["toLower"] = lower()
	m.Value["toUpper"] = upper()
	m.Value["count"] = count()
	m.Value["isAscii"] = isAscii()
	m.Value["isDecimal"] = isDecimal()
	m.Value["isDigit"] = isDigit()
	m.Value["isHexDigit"] = isHexDigit()
	m.Value["isLetter"] = isLetter()
	m.Value["isNumber"] = isNumber()
	m.Value["isSpace"] = isSpace()
	m.Value["codepoint"] = codepoint()
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
		return vida.String{Value: string(runes), Runes: runes}, nil
	}
}

func trim() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(vida.String); ok {
				if p, ok := args[1].(vida.String); ok {
					return vida.String{Value: strings.Trim(v.Value, p.Value)}, nil
				}
				return vida.String{Value: strings.Trim(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return vida.String{Value: strings.Trim(v.Value, " ")}, nil
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
					return vida.String{Value: strings.TrimLeft(v.Value, p.Value)}, nil
				}
				return vida.String{Value: strings.TrimLeft(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return vida.String{Value: strings.TrimLeft(v.Value, " ")}, nil
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
					return vida.String{Value: strings.TrimRight(v.Value, p.Value)}, nil
				}
				return vida.String{Value: strings.TrimRight(v.Value, " ")}, nil
			}
			return vida.NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(vida.String); ok {
				return vida.String{Value: strings.TrimRight(v.Value, " ")}, nil
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
					if vida.StringLength(v)*times > vida.MaxStringLen {
						return vida.NilValue, nil
					}
					return vida.String{Value: strings.Repeat(v.Value, int(times))}, nil
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
							return vida.String{Value: strings.Replace(s.Value, old.Value, nnew.Value, int(k))}, nil
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
						return vida.String{Value: strings.ReplaceAll(s.Value, old.Value, nnew.Value)}, nil
					}
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func center() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		l := len(args)
		if l == 2 {
			if str, ok := args[0].(vida.String); ok {
				if width, ok := args[1].(vida.Integer); ok {
					strlen := vida.StringLength(str)
					if width <= strlen {
						return str, nil
					}
					padding := width - strlen
					newString := str.Value
					sep := " "
					for i := vida.Integer(0); i < padding; i++ {
						if i%2 == 0 {
							newString = newString + sep
						} else {
							newString = sep + newString
						}
					}
					return vida.String{Value: newString}, nil
				}
			}
			return vida.NilValue, nil
		}
		if l > 2 {
			if str, ok := args[0].(vida.String); ok {
				if width, ok := args[1].(vida.Integer); ok {
					if sep, ok := args[2].(vida.String); ok {
						strlen := vida.StringLength(str)
						if width <= strlen {
							return str, nil
						}
						padding := width - strlen
						newString := str.Value
						for i := vida.Integer(0); i < padding; i++ {
							if i%2 == 0 {
								newString = newString + sep.Value
							} else {
								newString = sep.Value + newString
							}
						}
						return vida.String{Value: newString}, nil
					}
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func contains() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(vida.String); ok {
				if substr, ok := args[1].(vida.String); ok {
					return vida.Bool(strings.Contains(s.Value, substr.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func containsAny() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(vida.String); ok {
				if substr, ok := args[1].(vida.String); ok {
					return vida.Bool(strings.ContainsAny(s.Value, substr.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func index() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(vida.String); ok {
				if substr, ok := args[1].(vida.String); ok {
					return vida.Integer(strings.Index(s.Value, substr.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func join() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if xs, ok := args[0].(*vida.List); ok {
				if sep, ok := args[1].(vida.String); ok {
					var r []string
					for _, v := range xs.Value {
						r = append(r, v.String())
					}
					return vida.String{Value: strings.Join(r, sep.Value)}, nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func lower() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(vida.String); ok {
				return vida.String{Value: strings.ToLower(v.Value)}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func upper() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(vida.String); ok {
				return vida.String{Value: strings.ToUpper(v.Value)}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func count() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(vida.String); ok {
				if substr, ok := args[1].(vida.String); ok {
					return vida.Integer(strings.Count(s.Value, substr.Value)), nil
				}
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isAscii() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				c := s.Runes[0]
				return vida.Bool(0 <= c && c <= unicode.MaxASCII), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isDecimal() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				c := s.Runes[0]
				return vida.Bool('0' <= c && c <= '9'), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isDigit() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				return vida.Bool(unicode.IsDigit(s.Runes[0])), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isHexDigit() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				c := s.Runes[0]
				return vida.Bool('0' <= c && c <= '9' || 'a' <= (32|c) && (32|c) <= 'f'), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isLetter() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				c := s.Runes[0]
				return vida.Bool('a' <= (32|c) && (32|c) <= 'z' || c == '_' || c >= utf8.RuneSelf && unicode.IsLetter(c)), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isNumber() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				return vida.Bool(unicode.IsNumber(s.Runes[0])), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func isSpace() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				return vida.Bool(unicode.IsSpace(s.Runes[0])), nil
			}
			return vida.NilValue, nil
		}
		return vida.NilValue, nil
	}
}

func codepoint() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(vida.String); ok && vida.StringLength(s) == 1 {
				return vida.Integer(s.Runes[0]), nil
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
		xs[i] = vida.String{Value: slice[i]}
	}
	return &vida.List{Value: xs}
}
