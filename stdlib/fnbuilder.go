package stdlib

import "github.com/ever-eduardo/vida"

func FnRandInt(fn func(int64) int64) vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(vida.Integer); ok {
				if v > 0 {
					return vida.Integer(fn(int64(v))), nil
				}
				if v < 0 {
					return vida.Integer(fn(int64(-v))), nil
				}
				return vida.Integer(0), nil
			}
		}
		return vida.NilValue, nil
	}
}

func NormFloat(fn func() float64) vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		return vida.Float(fn()), nil
	}
}
