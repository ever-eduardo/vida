package stdlib

import (
	"math/rand/v2"

	"github.com/ever-eduardo/vida"
)

func generateRandom() vida.Value {
	m := &vida.Object{Value: make(map[string]vida.Value)}
	m.Value["int"] = FnRandInt(rand.Int64N)
	m.Value["normFloat"] = NormFloat(rand.NormFloat64)
	m.Value["float"] = NormFloat(rand.Float64)
	m.UpdateKeys()
	return m
}
