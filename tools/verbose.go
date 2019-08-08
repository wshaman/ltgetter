package tools

import "fmt"

func Verbose(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
