package libs

import (
	"math/rand/v2"

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
	m.Value["shuffled"] = randShuffled()
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

func randShuffled() vida.GFn {
	return func(args ...vida.Value) (vida.Value, error) {
		if len(args) > 0 {
			if v, ok := args[0].(*vida.List); ok {
				c := v.Clone().(*vida.List)
				rand.Shuffle(len(v.Value), func(i, j int) { c.Value[i], c.Value[j] = c.Value[j], c.Value[i] })
				return c, nil
			}
			if v, ok := args[0].(*vida.String); ok {
				if v.Runes == nil {
					v.Runes = []rune(v.Value)
				}
				l := len(v.Runes)
				r := make([]rune, l)
				copy(r, v.Runes)
				rand.Shuffle(l, func(i, j int) { r[i], r[j] = r[j], r[i] })
				return &vida.String{Value: string(r), Runes: r}, nil
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
