package stdlib

import "github.com/ever-eduardo/vida"

func LoadStdlib() map[string]func() vida.Value {
	l := make(map[string]func() vida.Value)
	l["rand"] = generateRandom
	l["math"] = generateMath
	l["text"] = generateText
	l["stringTo"] = generateStringTo
	return l
}
