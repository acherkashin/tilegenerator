package utils

import "os"

func FileExists(name *string) bool {
	if _, err := os.Stat(*name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir() && err != nil
}
