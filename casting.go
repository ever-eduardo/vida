package vida

import (
	"fmt"
	"strconv"
)

func loadFoundationCasting() Value {
	m := &Object{Value: make(map[string]Value, 6)}
	m.Value["toString"] = NativeFunction(castToString)
	m.Value["toInt"] = NativeFunction(castToInt)
	m.Value["toFloat"] = NativeFunction(castToFloat)
	m.Value["toBool"] = NativeFunction(castToBool)
	m.Value["toArray"] = NativeFunction(castToArray)
	m.Value["toObject"] = NativeFunction(castToObject)
	return m
}

func castToNumber(ctx *Context, args ...Value) (Value, error) {
	switch len(args) {
	case 2:
		inputStr, okInput := args[0].(*String)
		if !okInput {
			return argError("toNum", fmt.Sprintf("argument 1 must be a string, got %s", args[0].Type()))
		}
		typeOf, okType := args[1].(*String)
		if !okType {
			return argError("toNum", fmt.Sprintf("argument 2 (type) must be a string, got %s", args[1].Type()))
		}
		switch typeOf.Value {
		case integerT:
			i, e := strconv.ParseInt(inputStr.Value, 0, 64)
			if e == nil {
				return Integer(i), nil
			}
			return &VidaError{Message: &String{Value: e.Error()}}, nil
		case floatT:
			r, e := strconv.ParseFloat(inputStr.Value, 64)
			if e == nil {
				return Float(r), nil
			}
			return &VidaError{Message: &String{Value: e.Error()}}, nil
		default:
			return argError("toNum", fmt.Sprintf("argument 2 (type) must be %q or %q, got %q", integerT, floatT, typeOf.Value))
		}
	case 3:
		inputStr, okInput := args[0].(*String)
		if !okInput {
			return argError("toNum", fmt.Sprintf("argument 1 must be a string, got %s", args[0].Type()))
		}
		typeOf, okType := args[1].(*String)
		if !okType {
			return argError("toNum", fmt.Sprintf("argument 2 (type) must be a string, got %s", args[1].Type()))
		}
		base, okBase := args[2].(Integer)
		if !okBase {
			return argError("toNum", fmt.Sprintf("argument 3 (base) must be an integer, got %s", args[2].Type()))
		}
		if !(base == 0 || (2 <= base && base <= 36)) {
			return argError("toNum", fmt.Sprintf("base must be 0 or between 2 and 36, got %d", base))
		}
		switch typeOf.Value {
		case integerT:
			i, e := strconv.ParseInt(inputStr.Value, 0, int(base))
			if e == nil {
				return Integer(i), nil
			}
			return &VidaError{Message: &String{Value: e.Error()}}, nil
		case floatT:
			r, e := strconv.ParseFloat(inputStr.Value, 64)
			if e == nil {
				return Float(r), nil
			}
			return &VidaError{Message: &String{Value: e.Error()}}, nil
		default:
			return argError("toNum", fmt.Sprintf("argument 2 (type) must be %q or %q, got %q", integerT, floatT, typeOf.Value))
		}
	default:
		return argError("toNum", fmt.Sprintf("expected 2 arguments (string, type) or 3 arguments (string, type, base), got %d", len(args)))
	}
}

func castToString(ctx *Context, args ...Value) (Value, error) {
	if len(args) > 0 {
		return &String{Value: args[0].String()}, nil
	}
	return &String{Value: EmptyString}, nil
}

func castToInt(ctx *Context, args ...Value) (Value, error) {
	switch len(args) {
	case 1:
		switch v := args[0].(type) {
		case *String:
			i, e := strconv.ParseInt(v.Value, 0, 64)
			if e == nil {
				return Integer(i), nil
			}
			f, e := strconv.ParseFloat(v.Value, 64)
			if e == nil {
				return Integer(f), nil
			}
			return argError("toInt", fmt.Sprintf("cannot parse %q as a number", v.Value))
		case Integer:
			return v, nil
		case Bool:
			if v {
				return Integer(1), nil
			}
			return Integer(0), nil
		case Float:
			return Integer(int64(v)), nil
		case NilValue:
			return Integer(0), nil
		default:
			return argError("toInt", fmt.Sprintf("cannot convert %s to int", v.Type()))
		}
	case 2:
		v, ok := args[0].(*String)
		if !ok {
			return argError("toInt", fmt.Sprintf("a base can only be given when converting a string; argument 1 was %s", args[0].Type()))
		}
		base, okBase := args[1].(Integer)
		if !okBase {
			return argError("toInt", fmt.Sprintf("argument 2 (base) must be an integer, got %s", args[1].Type()))
		}
		if !(base == 0 || (2 <= base && base <= 36)) {
			return argError("toInt", fmt.Sprintf("base must be 0 or between 2 and 36, got %d", base))
		}
		i, e := strconv.ParseInt(v.Value, int(base), 64)
		if e == nil {
			return Integer(i), nil
		}
		return argError("toInt", fmt.Sprintf("cannot parse %q as a base-%d integer", v.Value, base))
	default:
		return argError("toInt", fmt.Sprintf("expected 1 or 2 arguments, got %d", len(args)))
	}
}

func castToFloat(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return Float(0), nil
	}
	switch v := args[0].(type) {
	case *String:
		r, e := strconv.ParseFloat(v.Value, 64)
		if e == nil {
			return Float(r), nil
		}
		return argError("toFloat", fmt.Sprintf("cannot parse %q as a number", v.Value))
	case Integer:
		return Float(v), nil
	case Float:
		return v, nil
	case NilValue:
		return Float(0), nil
	case Bool:
		if v {
			return Float(1), nil
		}
		return Float(0), nil
	default:
		return argError("toFloat", fmt.Sprintf("cannot convert %s to float", v.Type()))
	}
}

func castToBool(ctx *Context, args ...Value) (Value, error) {
	if len(args) > 0 {
		if v, ok := args[0].(*String); ok {
			switch v.Value {
			case "true":
				return True, nil
			case "false":
				return False, nil
			}
		}
		return args[0].Boolean(), nil
	}
	return False, nil
}

func castToArray(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return &Array{Value: make([]Value, 0)}, nil
	}
	switch v := args[0].(type) {
	case *Array:
		return v.Clone(), nil
	case *Bytes:
		a := make([]Value, len(v.Value))
		for i, b := range v.Value {
			a[i] = Integer(b)
		}
		return &Array{Value: a}, nil
	case *Object:
		pairs := make([]Value, len(v.Value))
		i := 0
		for k, val := range v.Value {
			pairs[i] = &Array{Value: []Value{&String{Value: k}, val}}
			i++
		}
		return &Array{Value: pairs}, nil
	case *Enum:
		pairs := make([]Value, len(v.Pairs))
		i := 0
		for k, val := range v.Pairs {
			pairs[i] = &Array{Value: []Value{&String{Value: k}, val}}
			i++
		}
		return &Array{Value: pairs}, nil
	case *VidaError:
		a := make([]Value, 1)
		a[0] = v.Message
		return &Array{Value: a}, nil
	default:
		return argError("toArray", fmt.Sprintf("cannot convert %s to array", v.Type()))
	}
}

func castToObject(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return &Object{Value: make(map[string]Value)}, nil
	}
	switch val := args[0].(type) {
	case *Object:
		return val.Clone(), nil
	case *Array:
		o := &Object{Value: make(map[string]Value, len(val.Value))}
		for i, v := range val.Value {
			o.Value[Integer(i).ObjectKey()] = v
		}
		return o, nil
	case *Bytes:
		o := &Object{Value: make(map[string]Value, len(val.Value))}
		for i, v := range val.Value {
			o.Value[Integer(i).ObjectKey()] = Integer(v)
		}
		return o, nil
	case *VidaError:
		o := &Object{Value: make(map[string]Value, 1)}
		o.Value[errorMessageFieldName] = val.Message
		return o, nil
	case *Enum:
		o := &Object{Value: make(map[string]Value, len(val.Pairs))}
		for k, v := range val.Pairs {
			o.Value[k] = v
		}
		return o, nil
	default:
		o := &Object{Value: make(map[string]Value, 4)}
		o.Value[DefaultValField] = val
		return o, nil
	}
}
