package compiler_test

import (
	"tim/forth/core/compiler"
	"fmt"
	"testing"
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
	err, result := compiler.Complete()

	if err != nil {
		t.Error("Got an error on complete")
		return
	}

	for key, value := range result {
		fmt.Println(key)
		fmt.Printf("%d\n", len(value))
	}

	if exp, found := result["refib"]; found {
		if len(exp) != 3 {
			t.Error("Base Expression wrong length")
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
