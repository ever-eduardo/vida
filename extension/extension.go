package extension

import "github.com/alkemist-17/vida"

var Success = vida.Bool(true)

func LoadExtensions() vida.LibsLoader {
	l := make(map[string]func() vida.Value)
	return l
}
