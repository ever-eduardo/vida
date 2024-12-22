package vida

import (
	"os"
)

const vidaFileExtension = ".vida"

func readModule(moduleName string) ([]byte, error) {
	if data, err := os.ReadFile(moduleName); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}
