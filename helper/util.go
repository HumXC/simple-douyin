package helper

import "os"

func IsFileExit(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
