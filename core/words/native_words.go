package words

import (
	"fmt"
	"tim/forth/core/support/stacks"
)

type UnderflowError struct{}

func (u UnderflowError) Error() string {
	return "Stack Underflow"
}
func NewUnderflowError() error {
	return UnderflowError{}
}

type InvalidArgument struct {
	message string
}

func (ia *InvalidArgument) Error() string {
	return fmt.Sprintf("Invalid argument: %s", ia.message)
}
func NewInvalidArgument(message string) error {
	return &InvalidArgument{
		message: message,
	}
}

func uinaryOperation(op func(int64) (int64, error)) func(*stacks.ForthStack) error {
	return func(stack *stacks.ForthStack) error {
		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		item := stack.Pop()

		result, err := op(item.ValueOf())
		if err != nil {
			return err
		}
		stack.Push(stacks.Number{Value: result})

		return nil
	}
}
func binaryOperation(op func(int64, int64) (int64, error)) func(*stacks.ForthStack) error {
	return func(stack *stacks.ForthStack) error {
		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		item1 := stack.Pop()

		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		item2 := stack.Pop()

		result, err := op(item1.ValueOf(), item2.ValueOf())
		if err != nil {
			return err
		}
		stack.Push(stacks.Number{Value: result})

		return nil
	}
}

func toStackBoolean(b bool) int64 {
	boolT := int64(0)
	boolF := int64(1)

	if b {
		return boolT
	}

	return boolF
}

func comparisonOperation(op func(int64, int64) bool) func(*stacks.ForthStack) error {
	return func(stack *stacks.ForthStack) error {
		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		item1 := stack.Pop()

		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		item2 := stack.Pop()

		result := op(item1.ValueOf(), item2.ValueOf())

		stackBoolean := toStackBoolean(result)

		stack.Push(item2)
		stack.Push(item1)
		stack.Push(stacks.Number{Value: stackBoolean})

		return nil
	}
}

func withItems(stack *stacks.ForthStack, n int, ifItems func([]stacks.ForthItem) error) error {
	items := make([]stacks.ForthItem, n)
	success := false
	index := 0
	for {
		if index >= n {
			success = true
			break
		}

		if stack.IsEmpty() {
			break
		}

		items[index] = stack.Pop()
		index = index + 1
	}

	if success {
		return ifItems(items)
	}

	return NewUnderflowError()
}

func NativeWords() map[string]func(*stacks.ForthStack) error {

	predefined := make(map[string]func(*stacks.ForthStack) error)
	predefined["dup"] = func(stack *stacks.ForthStack) error {
		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		toDuplicate := stack.Pop()
		stack.Push(toDuplicate)
		stack.Push(toDuplicate)

		return nil
	}
	predefined["drop"] = func(stack *stacks.ForthStack) error {
		if stack.IsEmpty() {
			return NewUnderflowError()
		}

		stack.Pop()

		return nil
	}
	predefined["print"] = func(stack *stacks.ForthStack) error {
		fmt.Printf("H -> %s <- T\n", stack.ToString())

		return nil
	}
	predefined["."] = func(stack *stacks.ForthStack) error {
		fmt.Println(stack.Peek().ToString())
		return nil
	}
	predefined["flip"] = func(stack *stacks.ForthStack) error {
		return withItems(stack, 2, func(items []stacks.ForthItem) error {
			v1 := items[0]
			v2 := items[1]

			stack.Push(v1)
			stack.Push(v2)

			return nil
		})
	}
	predefined["rotate"] = func(stack *stacks.ForthStack) error {
		return withItems(stack, 3, func(items []stacks.ForthItem) error {
			first := items[0]
			second := items[1]
			third := items[2]

			stack.Push(first)
			stack.Push(third)
			stack.Push(second)

			return nil
		})
	}
	predefined["+"] = binaryOperation(func(a int64, b int64) (int64, error) {
		return a + b, nil
	})
	predefined["*"] = binaryOperation(func(a int64, b int64) (int64, error) {
		return a * b, nil
	})
	predefined["-"] = binaryOperation(func(a int64, b int64) (int64, error) {
		return a - b, nil
	})

	predefined[">"] = comparisonOperation(func(a int64, b int64) bool {
		return (a > b)
	})
	predefined[">="] = comparisonOperation(func(a int64, b int64) bool {
		return (a >= b)
	})
	predefined["<"] = comparisonOperation(func(a int64, b int64) bool {
		return (a < b)
	})
	predefined["<="] = comparisonOperation(func(a int64, b int64) bool {
		return (a <= b)
	})
	predefined["=="] = comparisonOperation(func(a int64, b int64) bool {
		return (a == b)
	})
	predefined["!="] = comparisonOperation(func(a int64, b int64) bool {
		return (a != b)
	})
	predefined["branch"] = uinaryOperation(func(a int64) (int64, error) {
		if a == 1 {
			fmt.Println("True!!!")
			return a, nil
		}

		if a == 0 {
			fmt.Println("False!!!")
			return a, nil
		}

		return 0, NewInvalidArgument(fmt.Sprintf("The provided value [%d] is not acceptable as a boolean (must be 0 or 1)", a))
	})

	return predefined
}
