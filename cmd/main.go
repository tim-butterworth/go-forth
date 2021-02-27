package main

import (
	"fmt"
	"strings"
	"tim/forth/core"
	io "tim/forth/io/commandline"
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

	io.CommandLineSource(&handler{
		interpreter: core.NewForthInterpreter(),
	})
}
