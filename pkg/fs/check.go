package fs

import (
	"fmt"
	"os"
)

// MakeSure is used to make sure the path exists, or create it.
func MakeSure(path string) error {
	fi, err := os.Stat(path)
	if err == nil && !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("%s stat failed", path)
		}

		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("create directory(%s) failed: %v", path, err)
		}
	}

	return nil
}
