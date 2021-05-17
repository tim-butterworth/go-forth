package internal_test

import (
	"fmt"
	"testing"
	"tim/forth/core/support/stacks/internal"
)

func Test_BaseStackInitiallyEmpty(t *testing.T) {
	stack := internal.NewBackingStack()

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be initially empty")
}

func Test_IsNotEmpty_afterPushingAValue(t *testing.T) {
	stack := internal.NewBackingStack()

	stack.Push(int64(1))

	if stack.IsEmpty() {
		t.Error("Stack should not be empty if it has a value")
	}
}

func Test_IsEmpty_afterPushAndPopOfAValue(t *testing.T) {
	stack := internal.NewBackingStack()

	stack.Push(int64(1))
	stack.Pop(
		func(v int64) {},
		func() {
			t.Error("Should not have called the hasNoValue side")
		},
	)

	if stack.IsEmpty() {
		return
	}

	t.Error("Stack should be empty")
}

func Test_PusedValues_returnInReverseOrder(t *testing.T) {
	stack := internal.NewBackingStack()

	i := 0
	for {
		if i >= 10 {
			break
		}

		stack.Push(int64(i))

		i += 1
	}

	expectedResult := 9
	for {
		stack.Pop(
			func(v int64) {
				if v != int64(expectedResult) {
					message := fmt.Sprintf("Wrong Result expected %d  but got %d", expectedResult, v)
					t.Error(message)
				}
			},
			func() {
				t.Error("Should not be empty")
			},
		)

		if stack.IsEmpty() {
			break
		}
		expectedResult -= 1
	}
}

func Test_EmptyFnCalled_onEmptyStack(t *testing.T) {
	stack := internal.NewBackingStack()

	called := false
	stack.Pop(
		func(v int64) {
			message := fmt.Sprintf("Expected no invocation, instead called with [%d]", v)
			t.Error(message)
		},
		func() {
			called = true
		},
	)

	if called {
		return
	}

	t.Error("Should have invoked the no value function")
}

func Test_FillEmptyFillEmpty_workflow(t *testing.T) {
	stack := internal.NewBackingStack()

	fillLevel1 := 100
	fillStackUntil(stack, fillLevel1)

	if stack.IsEmpty() {
		t.Error("Should not be empty")
	}

	count := emptyStack(stack)

	if !stack.IsEmpty() {
		t.Error("Stack should be empty")
	}
	if count != fillLevel1 {
		t.Error(fmt.Sprintf("Expected to empty [%d] instead [%d]", fillLevel1, count))
	}

	fillLevel2 := 201
	fillStackUntil(stack, fillLevel2)

	count = emptyStack(stack)
	if count != fillLevel2 {
		t.Error(fmt.Sprintf("Expected to empty [%d] instead [%d]", fillLevel2, count))
	}
}

func fillStackUntil(stack internal.BackingStack, max int) {
	current := int64(0)
	max64 := int64(max)
	for {
		if current >= max64 {
			break
		}

		stack.Push(current)

		current += 1
	}
}

func emptyStack(stack internal.BackingStack) int {
	removed := 0
	for {
		if stack.IsEmpty() {
			break
		}
		stack.Pop(
			func(i int64) {},
			func() {},
		)

		removed += 1
	}
	return removed
}
