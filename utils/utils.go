package utils

import "os"

// 是否为开发环境
func IsDev() bool {
	args := os.Args
	for i := 0; i < len(args); i++ {
		if args[i] == "--dev" {
			return true
		}
	}
	return false
}
