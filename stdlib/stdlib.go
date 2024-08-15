package stdlib

import "github.com/ever-eduardo/vida"

func LoadLibs() map[string]func() vida.Value {
	l := make(map[string]func() vida.Value)
	l["rand"] = generateRandom
	return l
}
