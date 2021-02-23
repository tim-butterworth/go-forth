package io

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"tim/forth/core"
)

func Hi(commandHandler core.ForthCommandHandler) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("----------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("bye", text) == 0 {
			fmt.Println("Bye for now!")
			break
		}

		if strings.Compare("h", text) == 0 {
			fmt.Println("Help coming soon!")
		} else {
			fmt.Println(commandHandler.Execute(text))
		}
	}
}
