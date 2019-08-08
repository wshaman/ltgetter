package main

import (
	"fmt"

	"github.com/wshaman/ltntreader/command"
	"github.com/wshaman/ltntreader/tools"
)

func main() {
	fmt.Println("Hello from litnet parser")
	tools.OnErrPanic(command.Start())
}
