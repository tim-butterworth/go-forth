package main

import (
	"fmt"
	"strings"
	io "tim/forth/io/commandline"
	"tim/forth/core"
)

type handler struct {
	interpreter *core.ForthInterpreter
}
func (h handler) Execute(command string) string {
	splits := strings.Split(command, " ")
	for _, v := range splits {
		h.interpreter.Execute(v)
	}

	return "Consider it handled!"
}

func main() {
	fmt.Println("hi")

	io.Hi(&handler{
		interpreter: core.NewForthInterpreter(),
	})
}
