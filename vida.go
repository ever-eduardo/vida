package vida

import "fmt"

const v uint64 = 'v'
const i uint64 = 'i'
const d uint64 = 'd'
const a uint64 = 'a'

const major uint64 = 0
const minor uint64 = 3
const patch uint64 = 7
const inception uint64 = 24
const header uint64 = v<<56 | i<<48 | d<<40 | a<<32 | major<<24 | minor<<16 | patch<<8 | inception
const name = "Vida 🌱🐝🌻"

func Name() string {
	return name
}

func Version() string {
	return fmt.Sprintf("Version %v.%v.%v", major, minor, patch)
}

func About() string {
	return "\n\n\t" + name + "\n\t" + Version() + `
	
	
	What is Vida?

	Vida is a simple, interpreted computer language.
	The current implementation is written in Go.
	It has a minimal set of language constructs making it ergonomic
	and suitable for common tasks.
	The language can be extended by implementing the Value interface.
	The stdlib is implemented in such way, and you can add to the
	stdlib all you need for your projects.
	Contributions are welcome.
	
	Happy Vida coding!


	`
}
