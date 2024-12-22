package vida

import (
	"os"

	"github.com/ever-eduardo/vida/verror"
)

const vidaFileExtension = ".vida"

func readModule(moduleName string) ([]byte, error) {
	if data, err := os.ReadFile(moduleName); err == nil {
		return data, nil
	} else {
		return nil, verror.New(moduleName, err.Error(), verror.FileErrType, 0)
	}
}
