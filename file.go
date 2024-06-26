package vida

import (
	"os"
)

func ReadFile(filename string) ([]byte, error) {
	if data, err := os.ReadFile(filename); err == nil {
		return data, nil
	} else {
		return nil, err
	}
}
