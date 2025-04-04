package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func HashString(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateDirIfNotExist(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil
	}
	return nil
}
