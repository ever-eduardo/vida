package vida

import (
	"encoding/json"
	"fmt"
)

func loadFoundationJSON() Value {
	m := &Object{Value: make(map[string]Value, 6)}
	m.Value["stringify"] = NativeFunction(jsonValueToJsonString)
	m.Value["parse"] = NativeFunction(jsonParse)
	m.Value["encode"] = NativeFunction(jsonValueToJsonBytes)
	m.Value["decode"] = NativeFunction(jsonDecode)
	m.Value["isValid"] = NativeFunction(jsonIsValid)
	m.Value["pretty"] = NativeFunction(jsonPretty)
	return m
}

func jsonValueToJsonString(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.stringify", "expected a value argument")
	}
	if b, err := json.Marshal(args[0]); err == nil {
		return &String{Value: string(b)}, nil
	} else {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
}

func jsonValueToJsonBytes(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.encode", "expected a value argument")
	}
	if b, err := json.Marshal(args[0]); err == nil {
		return &Bytes{Value: b}, nil
	} else {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
}

func jsonParse(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.parse", "expected a string or bytes argument")
	}
	valid, _ := jsonIsValid(ctx, args[0])
	switch t := valid.(type) {
	case NilValue:
		return &VidaError{Message: &String{Value: ErrInvalidJSON.Error()}}, nil
	case Bool:
		if !t {
			return &VidaError{Message: &String{Value: ErrInvalidJSON.Error()}}, nil
		}
	}
	switch t := args[0].(type) {
	case *String:
		var value any
		if err := json.Unmarshal([]byte(t.Value), &value); err == nil {
			switch v := value.(type) {
			case nil:
				return Nil, nil
			case bool:
				return Bool(v), nil
			case string:
				return &String{Value: v}, nil
			case float64:
				return Float(v), nil
			case map[string]any:
				return parseObject(ctx, v), nil
			case []any:
				return parseArray(ctx, v), nil
			}
			return Nil, nil
		} else {
			return &VidaError{Message: &String{Value: err.Error()}}, nil
		}
	case *Bytes:
		var value any
		if err := json.Unmarshal(t.Value, &value); err == nil {
			switch v := value.(type) {
			case nil:
				return Nil, nil
			case bool:
				return Bool(v), nil
			case string:
				return &String{Value: v}, nil
			case float64:
				return Float(v), nil
			case map[string]any:
				return parseObject(ctx, v), nil
			case []any:
				return parseArray(ctx, v), nil
			}
			return Nil, nil
		} else {
			return &VidaError{Message: &String{Value: err.Error()}}, nil
		}
	default:
		return argError("json.parse", fmt.Sprintf("expected a string or bytes argument, got %s", args[0].Type()))
	}
}

func jsonDecode(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.decode", "expected a bytes argument")
	}
	valid, _ := jsonIsValid(ctx, args[0])
	switch t := valid.(type) {
	case NilValue:
		return &VidaError{Message: &String{Value: ErrInvalidJSON.Error()}}, nil
	case Bool:
		if !t {
			return &VidaError{Message: &String{Value: ErrInvalidJSON.Error()}}, nil
		}
	}
	switch t := args[0].(type) {
	case *Bytes:
		var value any
		if err := json.Unmarshal(t.Value, &value); err == nil {
			switch v := value.(type) {
			case nil:
				return Nil, nil
			case bool:
				return Bool(v), nil
			case string:
				return &String{Value: v}, nil
			case float64:
				return Float(v), nil
			case map[string]any:
				return parseObject(ctx, v), nil
			case []any:
				return parseArray(ctx, v), nil
			}
			return Nil, nil
		} else {
			return &VidaError{Message: &String{Value: err.Error()}}, nil
		}
	default:
		return argError("json.decode", fmt.Sprintf("expected a bytes argument, got %s", args[0].Type()))
	}
}

func jsonIsValid(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.isValid", "expected a string or bytes argument")
	}
	switch t := args[0].(type) {
	case *String:
		return Bool(json.Valid([]byte(t.Value))), nil
	case *Bytes:
		return Bool(json.Valid(t.Value)), nil
	default:
		return argError("json.isValid", fmt.Sprintf("expected a string or bytes argument, got %s", t.Type()))
	}
}

func jsonPretty(ctx *Context, args ...Value) (Value, error) {
	if len(args) == 0 {
		return argError("json.pretty", "expected a value argument")
	}
	if b, err := json.MarshalIndent(args[0], EmptyString, "  "); err == nil {
		return &String{Value: string(b)}, nil
	} else {
		return &VidaError{Message: &String{Value: err.Error()}}, nil
	}
}

func parseObject(ctx *Context, input map[string]any) *Object {
	o := &Object{Value: make(map[string]Value, len(input))}
	for kk, vv := range input {
		switch tt := vv.(type) {
		case nil:
			o.Value[kk] = Nil
		case bool:
			o.Value[kk] = Bool(tt)
		case string:
			o.Value[kk] = &String{Value: tt}
		case float64:
			o.Value[kk] = Float(tt)
		case map[string]any:
			o.Value[kk] = parseObject(ctx, tt)
		case []any:
			o.Value[kk] = parseArray(ctx, tt)
		}
	}
	return o
}

func parseArray(ctx *Context, input []any) *Array {
	A := &Array{Value: make([]Value, len(input))}
	for ii, vv := range input {
		switch tt := vv.(type) {
		case nil:
			A.Value[ii] = Nil
		case bool:
			A.Value[ii] = Bool(tt)
		case string:
			A.Value[ii] = &String{Value: tt}
		case float64:
			A.Value[ii] = Float(tt)
		case map[string]any:
			A.Value[ii] = parseObject(ctx, tt)
		case []any:
			A.Value[ii] = parseArray(ctx, tt)
		}
	}
	return A
}
