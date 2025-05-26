package vida

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func loadFoundationText() Value {
	m := &Object{Value: make(map[string]Value)}
	m.Value["hasPrefix"] = hasPrefix()
	m.Value["hasSuffix"] = hasSuffix()
	m.Value["fromCodePoint"] = fromCodepoint()
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
	m.Value["codePoint"] = codepoint()
	m.Value["byteslen"] = byteslen()
	m.UpdateKeys()
	return m
}

func hasPrefix() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return Bool(strings.HasPrefix(v.Value, p.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func hasSuffix() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return Bool(strings.HasSuffix(v.Value, p.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func fromCodepoint() GFn {
	return func(args ...Value) (Value, error) {
		runes := make([]rune, 0)
		for _, a := range args {
			if v, ok := a.(Integer); ok && v > 0 {
				runes = append(runes, int32(v))
			}
		}
		return &String{Value: string(runes), Runes: runes}, nil
	}
}

func trim() GFn {
	return func(args ...Value) (Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return &String{Value: strings.Trim(v.Value, p.Value)}, nil
				}
				return &String{Value: strings.Trim(v.Value, " ")}, nil
			}
			return NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(*String); ok {
				return &String{Value: strings.Trim(v.Value, " ")}, nil
			}
		}
		return NilValue, nil
	}
}

func trimLeft() GFn {
	return func(args ...Value) (Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return &String{Value: strings.TrimLeft(v.Value, p.Value)}, nil
				}
				return &String{Value: strings.TrimLeft(v.Value, " ")}, nil
			}
			return NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(*String); ok {
				return &String{Value: strings.TrimLeft(v.Value, " ")}, nil
			}
		}
		return NilValue, nil
	}
}

func trimRight() GFn {
	return func(args ...Value) (Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return &String{Value: strings.TrimRight(v.Value, p.Value)}, nil
				}
				return &String{Value: strings.TrimRight(v.Value, " ")}, nil
			}
			return NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(*String); ok {
				return &String{Value: strings.TrimRight(v.Value, " ")}, nil
			}
		}
		return NilValue, nil
	}
}

func split() GFn {
	return func(args ...Value) (Value, error) {
		l := len(args)
		if l > 1 {
			if v, ok := args[0].(*String); ok {
				if p, ok := args[1].(*String); ok {
					return stringSliceToList(strings.Split(v.Value, p.Value)), nil
				}
				return stringSliceToList(strings.Split(v.Value, "")), nil
			}
			return NilValue, nil
		}
		if l == 1 {
			if v, ok := args[0].(*String); ok {
				return stringSliceToList(strings.Split(v.Value, "")), nil
			}
		}
		return NilValue, nil
	}
}

func fields() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*String); ok {
				return stringSliceToList(strings.Fields(v.Value)), nil
			}
		}
		return NilValue, nil
	}
}

func repeat() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) >= 2 {
			if v, ok := args[0].(*String); ok {
				if times, ok := args[1].(Integer); ok && times >= 0 {
					if StringLength(v)*times > MaxStringSize {
						return NilValue, nil
					}
					return &String{Value: strings.Repeat(v.Value, int(times))}, nil
				}
				return NilValue, nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func replace() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 3 {
			if s, ok := args[0].(*String); ok {
				if old, ok := args[1].(*String); ok {
					if nnew, ok := args[2].(*String); ok {
						if k, ok := args[3].(Integer); ok {
							return &String{Value: strings.Replace(s.Value, old.Value, nnew.Value, int(k))}, nil
						}
					}
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func replaceAll() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 2 {
			if s, ok := args[0].(*String); ok {
				if old, ok := args[1].(*String); ok {
					if nnew, ok := args[2].(*String); ok {
						return &String{Value: strings.ReplaceAll(s.Value, old.Value, nnew.Value)}, nil
					}
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func center() GFn {
	return func(args ...Value) (Value, error) {
		l := len(args)
		if l == 2 {
			if str, ok := args[0].(*String); ok {
				if width, ok := args[1].(Integer); ok {
					strlen := StringLength(str)
					if width <= strlen {
						return str, nil
					}
					padding := width - strlen
					newString := str.Value
					sep := " "
					for i := Integer(0); i < padding; i++ {
						if i%2 == 0 {
							newString = newString + sep
						} else {
							newString = sep + newString
						}
					}
					return &String{Value: newString}, nil
				}
			}
			return NilValue, nil
		}
		if l > 2 {
			if str, ok := args[0].(*String); ok {
				if width, ok := args[1].(Integer); ok {
					if sep, ok := args[2].(*String); ok {
						strlen := StringLength(str)
						if width <= strlen {
							return str, nil
						}
						padding := width - strlen
						newString := str.Value
						for i := Integer(0); i < padding; i++ {
							if i%2 == 0 {
								newString = newString + sep.Value
							} else {
								newString = sep.Value + newString
							}
						}
						return &String{Value: newString}, nil
					}
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func contains() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(*String); ok {
				if substr, ok := args[1].(*String); ok {
					return Bool(strings.Contains(s.Value, substr.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func containsAny() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(*String); ok {
				if substr, ok := args[1].(*String); ok {
					return Bool(strings.ContainsAny(s.Value, substr.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func index() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(*String); ok {
				if substr, ok := args[1].(*String); ok {
					return Integer(strings.Index(s.Value, substr.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func join() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if xs, ok := args[0].(*List); ok {
				if sep, ok := args[1].(*String); ok {
					var r []string
					for _, v := range xs.Value {
						r = append(r, v.String())
					}
					return &String{Value: strings.Join(r, sep.Value)}, nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func lower() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*String); ok {
				return &String{Value: strings.ToLower(v.Value)}, nil
			}
		}
		return NilValue, nil
	}
}

func upper() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*String); ok {
				return &String{Value: strings.ToUpper(v.Value)}, nil
			}
		}
		return NilValue, nil
	}
}

func count() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 1 {
			if s, ok := args[0].(*String); ok {
				if substr, ok := args[1].(*String); ok {
					return Integer(strings.Count(s.Value, substr.Value)), nil
				}
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isAscii() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				c := s.Runes[0]
				return Bool(0 <= c && c <= unicode.MaxASCII), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isDecimal() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				c := s.Runes[0]
				return Bool('0' <= c && c <= '9'), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isDigit() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				return Bool(unicode.IsDigit(s.Runes[0])), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isHexDigit() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				c := s.Runes[0]
				return Bool('0' <= c && c <= '9' || 'a' <= (32|c) && (32|c) <= 'f'), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isLetter() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				c := s.Runes[0]
				return Bool('a' <= (32|c) && (32|c) <= 'z' || c == '_' || c >= utf8.RuneSelf && unicode.IsLetter(c)), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isNumber() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				return Bool(unicode.IsNumber(s.Runes[0])), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func isSpace() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				return Bool(unicode.IsSpace(s.Runes[0])), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func codepoint() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if s, ok := args[0].(*String); ok && StringLength(s) == 1 {
				return Integer(s.Runes[0]), nil
			}
			return NilValue, nil
		}
		return NilValue, nil
	}
}

func byteslen() GFn {
	return func(args ...Value) (Value, error) {
		if len(args) > 0 {
			if val, ok := args[0].(*String); ok {
				return Integer(len(val.Value)), nil
			}
		}
		return NilValue, nil
	}
}

func stringSliceToList(slice []string) Value {
	l := len(slice)
	xs := make([]Value, l)
	for i := 0; i < l; i++ {
		xs[i] = &String{Value: slice[i]}
	}
	return &List{Value: xs}
}
