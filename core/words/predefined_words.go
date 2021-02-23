package words

func PredefinedWords() map[string][]string {
	predefined := make(map[string][]string)

	predefined["square"] = []string{ "*", "dup" }
	predefined["fib"] = []string{ "+", "rotate", "flip", "dup"}
	predefined["fib-10"] = []string{".", "fib", "fib", "fib", "fib", "fib", "fib", "fib", "fib", "fib", "fib", "1", "1"}

	return predefined
}

