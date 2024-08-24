package stdlib

import (
	"math/rand/v2"
	"strings"

	"github.com/ever-eduardo/vida"
)

func generateRandom() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["nextI"] = randNextI()
	m.Value["nextU32"] = randNextU32()
	m.Value["nextF"] = randNextF(rand.Float64)
	m.Value["norm"] = randNextF(rand.NormFloat64)
	m.Value["exp"] = randNextF(rand.ExpFloat64)
	m.Value["perm"] = randPerm()
	m.Value["shuffle"] = randShuffle()
	m.UpdateKeys()
	return m
}

func randNextI() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) == 0 {
			return vida.Integer(rand.Int64()), nil
		}
		if len(args) > 0 {
			if v, ok := args[0].(vida.Integer); ok {
				if v > 0 {
					return vida.Integer(rand.Int64N(int64(v))), nil
				}
				if v < 0 {
					return vida.Integer(rand.Int64N(int64(-v))), nil
				}
				return vida.Integer(0), nil
			} else {
				return vida.Integer(rand.Int64()), nil
			}
		}
		return vida.NilValue, nil
	}
}

func randNextF(fn func() float64) vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		return vida.Float(fn()), nil
	}
}

func randPerm() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(vida.Integer); ok {
				if 0 <= v && v <= 0x7FFF_FFFF {
					xs := make([]vida.Value, v)
					for i := range xs {
						xs[i] = vida.Integer(i)
					}
					rand.Shuffle(int(v), func(i, j int) { xs[i], xs[j] = xs[j], xs[i] })
					return &vida.List{Value: xs}, nil
				}
			}
		}
		return vida.NilValue, nil
	}
}

func randShuffle() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.List); ok {
				rand.Shuffle(len(v.Value), func(i, j int) { v.Value[i], v.Value[j] = v.Value[j], v.Value[i] })
				return v, nil
			}
			if v, ok := args[0].(vida.String); ok {
				s := strings.Split(v.Value, "")
				rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
				return vida.String{Value: strings.Join(s, "")}, nil
			}
		}
		return vida.NilValue, nil
	}
}

func randNextU32() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) == 0 {
			return vida.Integer(rand.Int32()), nil
		}
		if len(args) > 0 {
			if v, ok := args[0].(vida.Integer); ok {
				if v > 0 {
					return vida.Integer(uint32(rand.Int64N(int64(v)))), nil
				}
				if v < 0 {
					return vida.Integer(uint32(rand.Int64N(int64(-v)))), nil
				}
				return vida.Integer(0), nil
			} else {
				return vida.Integer(rand.Int32()), nil
			}
		}
		return vida.NilValue, nil
	}
}
