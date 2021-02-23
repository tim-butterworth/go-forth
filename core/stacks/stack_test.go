package stacks_test

import (
	"fmt"
	"testing"
	"tim/forth/core/stacks"
)

func repeat(times int, op func()) {
	n := times
	for {
		if n <= 0 {
			break
		}

		op()
		n = n - 1
	}
}

func Test_StackStartsEmpty(t *testing.T) {
	stack := stacks.NewStack()

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be initially empty")
}

func Test_StackIsNotEmptyAfterPushing(t *testing.T) {
	stack := stacks.NewStack()

	stack.Push(stacks.Empty{})

	if stack.IsEmpty() {
		t.Error("Stack should not be empty")
	}
}

func Test_StackIsEmptyAfter_Push_Pop(t *testing.T) {
	stack := stacks.NewStack()

	stack.Push(stacks.Empty{})
	stack.Pop()

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be empty")
}

func Test_StackIsEmtpyAfter_nPushes_and_nPops(t *testing.T) {
	stack := stacks.NewStack()

	repeat(100, func() {
		stack.Push(stacks.Empty{})
	})

	repeat(99, func() {
		stack.Pop()

		if stack.IsEmpty() {
			t.Error("Stack should not be empty")
		}
	})

	stack.Pop()

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be empty")
}

func Test_StackIsEmptyAfter_Pop(t *testing.T) {
	stack := stacks.NewStack()

	stack.Pop()

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be empty")
}

func Test_StackToString(t *testing.T) {
	stack := stacks.NewStack()

	n := 0
	for {
		if n >= 10 {
			break
		}
		stack.Push(stacks.Number{Value: int64(n)})

		n = n + 1
	}

	printResult := stack.ToString()

	fmt.Println(printResult)

	if stack.IsEmpty() {
		t.Error("Stack should not be empty")
	}
}
