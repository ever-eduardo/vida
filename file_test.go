package vida

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	const source = "./scripts/hello.vida"
	if _, err := ReadModule(source); err != nil {
		t.Fatalf("Failed to read the source file %q", source)
	}
}
func TestReadFileWithNoFile(t *testing.T) {
	const source = "asudhoiusdhfoaidsf.vida"
	if _, err := ReadModule(source); err == nil {
		t.Fatalf("Failed to read the source file %q", source)
	}
}
