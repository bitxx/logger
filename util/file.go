package util

import (
	"os"
)

// PathExist 判断目录是否存在
func PathExist(addr string) bool {
	s, err := os.Stat(addr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}
