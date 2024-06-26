package main

import (
	"fmt"
	"hash/maphash"

	"github.com/ever-eduardo/vida"
)

func main() {
	println(vida.Name(), vida.Version())
	var h maphash.Hash
	h.WriteString("Hello, World!")
	fmt.Println(h.Sum64())
}
