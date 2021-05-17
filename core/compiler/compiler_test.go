package compiler_test

import (
	"fmt"
	"testing"
	"tim/forth/core/compiler"
	"tim/forth/core/support/stacks"
)

var ifS string = "if"
var elseS string = "else"
var thenS string = "then"

func Test_SupportsFullExpression(t *testing.T) {
	compiler := compiler.NewCompiler(&testIdProvider{
		current: 0,
	})

	command := []string{
		"refib",
		"0",
		">",
		ifS,
		"1",
		"flip",
		"-",
		"rotate",
		"fib",
		"rotate",
		"rotate",
		"refib",
		elseS,
		"drop",
		"print",
		thenS,
	}

	for _, word := range command {
		compiler.PushWord(word)
	}
	result, err := compiler.Complete()

	if err != nil {
		t.Error("Got an error on complete")
		return
	}

	if len(result) != 2 {
		t.Error(fmt.Sprintf("Expected %d entries in the function results, instead got %d", 2, len(result)))
	}
	for key := range result {
		fmt.Println(key)
	}

	forthStack := stacks.NewStack()
	executionStack := stacks.NewStringStack()
	if refib, found := result["refib"]; found {
		refib(forthStack, executionStack)

		fmt.Printf("forthStack -> %s <-\n", forthStack.ToString())
		fmt.Printf("executionStack -> %s <- \n", executionStack.ToString())

		if executionStack.Size() != 3 {
			t.Error(fmt.Sprintf("There should be 3 items in the execution stack, instead found [%d]", executionStack.Size()))
			return
		}

		firstEntry := executionStack.Pop()
		if firstEntry != "0" {
			t.Error(fmt.Sprintf("Expected [%s] but got [%s]", "0", firstEntry))
		}

		secondEntry := executionStack.Pop()
		if secondEntry != ">" {
			t.Error(fmt.Sprintf("Expected [%s] but got [%s]", ">", secondEntry))
		}
	} else {
		t.Error("There should have been a 'refib' entry")
	}
}

func Test_SupportsNestedIfs(t *testing.T) {

}

func Test_missingThen_error(t *testing.T) {
	compiler := compiler.NewCompiler(&testIdProvider{
		current: 0,
	})

	command := []string{
		"blah",
		"0",
		">",
		ifS,
		"1",
		elseS,
	}

	for _, word := range command {
		compiler.PushWord(word)
	}
	err, _ := compiler.Complete()

	if err == nil {
		t.Error("Should have gotten an error when then is missing")
	}
}

func Test_missingIf_hasElseThen_error(t *testing.T) {
	compiler := compiler.NewCompiler(&testIdProvider{
		current: 0,
	})

	command := []string{
		elseS,
		"1",
		thenS,
	}

	for _, word := range command {
		compiler.PushWord(word)
	}
	err, _ := compiler.Complete()

	if err == nil {
		t.Error("Should have gotten an error")
	}
}

func Test_missingIf_hasThen_error(t *testing.T) {
	compiler := compiler.NewCompiler(&testIdProvider{
		current: 0,
	})

	command := []string{
		"1",
		thenS,
	}

	for _, word := range command {
		compiler.PushWord(word)
	}
	err, _ := compiler.Complete()

	if err == nil {
		t.Error("Should have gotten an error")
	}
}

type testIdProvider struct {
	current int64
}

func (idProvider *testIdProvider) NextId() string {
	current := idProvider.current
	id := fmt.Sprintf("id_%d", current)

	idProvider.current = current + 1

	return id
}
