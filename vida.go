package vida

import "fmt"

const major = 0
const minor = 1
const patch = 0
const inception = 24
const name = "Vida 🌱 Programming Language"

func Name() string {
	return name
}

func Version() string {
	return fmt.Sprintf("Version %v.%v.%v", major, minor, patch)
}
