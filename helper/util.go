package helper

import (
	"math/rand"
	"os"
)

func IsFileExit(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// 从切片里随机选择一个元素
func PickOne[T any](list []T) T {
	result := *new(T)
	if len(list) == 0 {
		return result
	}
	i := rand.Intn(len(list) - 1)
	result = list[i]
	return result
}
