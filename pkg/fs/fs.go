package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// TravelDirectory travel directory and returns all files in the directory.
func TravelDirectory(root string, recursive bool) ([]string, error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("failed to travelDirectory: %v", err)
	}

	out := make([]string, 0, len(files))
	for _, fp := range files {
		if !fp.IsDir() || !recursive {
			// not a directory or not recursive
			out = append(out, filepath.Join(root, fp.Name()))
			continue
		}

		subFiles, err := TravelDirectory(filepath.Join(root, fp.Name()), recursive)
		if err != nil {
			return nil, err
		}
		out = append(out, subFiles...)
	}

	return out, nil
}
