package lib

import "github.com/ever-eduardo/vida"

var Success = &vida.String{Value: string(vida.Success)}

func Loadlibs() vida.LibsLoader {
	l := make(map[string]func() vida.Value)
	l["rand"] = generateRandom
	l["math"] = generateMath
	l["text"] = generateText
	l["io"] = generateIO
	l["os"] = generateOS
	return l
}
