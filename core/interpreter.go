package core

import (
	"fmt"
	"strconv"
	"tim/forth/core/stacks"
	"tim/forth/core/words"
)

type ForthInterpreter struct {
	stack              *stacks.ForthStack
	newWordAccumulator *newWordAccumulator
	words              map[string][]string
	nativeWords        map[string]func(*stacks.ForthStack) error
	handler            func(*ForthInterpreter, string)
}

func processCommand(i *ForthInterpreter, executionStack stacks.StringStack) {
	for {
		if executionStack.IsEmpty() {
			break
		}

		command := executionStack.Pop()

		num, err := strconv.ParseInt(command, 10, 64)
		if err == nil {
			i.stack.Push(stacks.Number{Value: num})
			continue
		}

		maybeWord := i.words[command]
		if maybeWord != nil {
			for _, w := range maybeWord {
				fmt.Println(w)
				executionStack.Push(w)
			}
			continue
		}

		maybeFun := i.nativeWords[command]
		if maybeFun == nil {
			break //should return an error here
		}

		err = maybeFun(i.stack)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func executeCommand(i *ForthInterpreter, s string) {
	executionStack := stacks.NewStringStack()
	executionStack.Push(s)

	processCommand(i, executionStack)
}

func saveNewWord(i *ForthInterpreter, _ string) {
	accumulator := i.newWordAccumulator
	newWordLabel := accumulator.label.value()

	body := make([]string, accumulator.body.Size())
	recordedBody := accumulator.body
	index := 0
	for {
		if recordedBody.IsEmpty() {
			break
		}

		body[index] = recordedBody.Pop()
		index = index + 1
	}

	fmt.Printf("Saving the new word! %s\n", newWordLabel)
	i.words[newWordLabel] = body

	i.handler = executeCommand
	i.newWordAccumulator = NewWordAccumulator()
}

func recordNewWord(i *ForthInterpreter, s string) {
	fmt.Println("Recording!!!")
	if s != ":" {
		i.newWordAccumulator.insert(s)
	}
}

func (i *ForthInterpreter) Execute(s string) {
	if s == ":" {
		i.handler = recordNewWord
	} else if s == ";" {
		i.handler = saveNewWord
	}

	i.handler(i, s)
}

type Label interface {
	isEmpty() bool
	value() string
}
type EmptyLabel struct {
}

func (e EmptyLabel) isEmpty() bool {
	return true
}
func (e EmptyLabel) value() string {
	return ""
}

type PopulatedLabel struct {
	label string
}

func (p PopulatedLabel) isEmpty() bool {
	return false
}
func (p PopulatedLabel) value() string {
	return p.label
}

func NewPopulatedLabel(s string) Label {
	return &PopulatedLabel{label: s}
}

type newWordAccumulator struct {
	label Label
	body  stacks.StringStack
}

func (a *newWordAccumulator) insert(s string) {
	if a.label.isEmpty() {
		a.label = NewPopulatedLabel(s)
	} else {
		a.body.Push(s)
	}
}

func NewWordAccumulator() *newWordAccumulator {
	return &newWordAccumulator{
		label: EmptyLabel{},
		body:  stacks.NewStringStack(),
	}
}

func NewForthInterpreter() *ForthInterpreter {
	return &ForthInterpreter{
		stack:              stacks.NewStack(),
		nativeWords:        words.NativeWords(),
		words:              words.PredefinedWords(),
		newWordAccumulator: NewWordAccumulator(),
		handler:            executeCommand,
	}
}
