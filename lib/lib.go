package lib

import "github.com/alkemist-17/vida"

var Success = &vida.String{Value: string(vida.Success)}

func Loadlibs() vida.LibsLoader {
	l := make(map[string]func() vida.Value)
	l["rand"] = generateRandom
	l["io"] = generateIO
	return l
}
