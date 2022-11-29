package gots

import (
	"path/filepath"
)

func GetAbsPath(filename string) (string, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	return abs, nil
}
