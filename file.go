package vida

import (
	"os"

	"github.com/alkemist-17/vida/verror"
)

const vidaFileExtension = ".vida"

func readScript(scriptName string) ([]byte, error) {
	if data, err := os.ReadFile(scriptName); err == nil {
		return data, nil
	} else {
		return nil, verror.New(scriptName, err.Error(), verror.FileErrType, 0)
	}
}
