package vida

import "fmt"

const v uint64 = 'v'
const i uint64 = 'i'
const d uint64 = 'd'
const a uint64 = 'a'

const major uint64 = 0
const minor uint64 = 2
const patch uint64 = 5
const inception uint64 = 24
const header uint64 = v<<56 | i<<48 | d<<40 | a<<32 | major<<24 | minor<<16 | patch<<8 | inception
const name = "Vida 🌱🌻🐝 Programming Language"

func Name() string {
	return name
}

func Version() string {
	return fmt.Sprintf("Version %v.%v.%v", major, minor, patch)
}
