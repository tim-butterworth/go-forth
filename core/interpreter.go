package core

import (
	"fmt"
	"strconv"
	"tim/forth/core/compiler"
	"tim/forth/core/support/stacks"
	"tim/forth/core/words"

	"github.com/google/uuid"
)

type ForthInterpreter struct {
	stack              *stacks.ForthStack
	newWordAccumulator *newWordAccumulator
	words              map[string]func(*stacks.ForthStack, stacks.StringStack) error
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

		fun, found := i.words[command]
		if !found {
			fmt.Printf("Word -> [%s] is not defined (the stack should probably be dumped in this case)\n", command)
			fmt.Println("Available commands are:")
			for key := range i.words {
				fmt.Println(key)
			}
			break
		}

		err = fun(i.stack, executionStack)
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

type wordEntry struct {
	key  string
	body []string
}

func (entry *wordEntry) getKey() string {
	return entry.key
}
func (entry *wordEntry) getBody() []string {
	return entry.body
}

type uuidProvider struct{}

func (idProvider *uuidProvider) NextId() string {
	id, err := uuid.NewUUID()
	if err != nil {
		return ""
	}

	return id.String()
}

func endRecording(i *ForthInterpreter, _ string) {
	compiler := compiler.NewCompiler(&uuidProvider{})

	accumulator := i.newWordAccumulator
	i.newWordAccumulator.label.value()
	compiler.PushWord(accumulator.label.value())

	for _, w := range accumulator.body[0:accumulator.wordCount] {
		compiler.PushWord(w)
	}

	words, err := compiler.Complete()
	if err != nil {
		fmt.Println("Failed to process some stuff :(")
	} else {
		for label, body := range words {
			i.words[label] = body
		}
	}

	i.handler = executeCommand
	i.newWordAccumulator = NewWordAccumulator()
}

func record(i *ForthInterpreter, s string) {
	i.newWordAccumulator.insert(s)
}
func startRecording(i *ForthInterpreter, _ string) {
	fmt.Println("Recording!!!")
	i.handler = record
}

func (i *ForthInterpreter) Execute(s string) {
	if s == ":" {
		i.handler = startRecording
	} else if s == ";" {
		i.handler = endRecording
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
	label     Label
	body      []string
	wordCount int32
}

func (a *newWordAccumulator) insert(s string) {
	if a.label.isEmpty() {
		a.label = NewPopulatedLabel(s)
	} else {
		a.body[a.wordCount] = s
		a.wordCount = a.wordCount + 1
	}
}

func NewWordAccumulator() *newWordAccumulator {
	return &newWordAccumulator{
		label:     EmptyLabel{},
		body:      make([]string, 100),
		wordCount: 0,
	}
}

func wrapNative(name string, fun func(*stacks.ForthStack) error) func(*stacks.ForthStack, stacks.StringStack) error {
	return func(forthStack *stacks.ForthStack, executionStack stacks.StringStack) error {
		fmt.Printf("Calling through to native function, [%s]\n", name)
		return fun(forthStack)
	}
}

func wrapPredefined(name string, body []string) func(*stacks.ForthStack, stacks.StringStack) error {
	return func(forthStack *stacks.ForthStack, executionStack stacks.StringStack) error {
		for _, w := range body {
			fmt.Println(w)
			executionStack.Push(w)
		}
		return nil
	}
}

func NewForthInterpreter() *ForthInterpreter {
	nativeWords := words.NativeWords()
	predefinedWords := words.PredefinedWords()

	words := make(map[string]func(*stacks.ForthStack, stacks.StringStack) error)
	for key, value := range nativeWords {
		words[key] = wrapNative(key, value)
	}

	for key, body := range predefinedWords {
		words[key] = wrapPredefined(key, body)
	}

	return &ForthInterpreter{
		stack:              stacks.NewStack(),
		words:              words,
		newWordAccumulator: NewWordAccumulator(),
		handler:            executeCommand,
	}
}
