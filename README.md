# This is a small, limited implementation of the programming language forth

### Prerequisites

- Golang (15 or greater, have not tested with anything older than 15)

### Commands

`repl` - from the root directory `go run cmd/main.go`

in the repl there are a few built in commands:

`bye` - exits the repl
`print` - prints the contents of the stack
`.` - shows the first value of the stack without popping it (peek)
`drop` - drops the top value of the stack
`dup` - duplicates the top value of the stack

There are more 'native' built in functions (functions implemented in go, not in forth) that can be found by checking the map in `core/words/native_words.go`

There are also predefined forth functions like:

`fib` - sums the top two numbers of the stack and leaves the stack in such a state that fib can be called again
`fib-10` - calculates some fibinacci numbers (populates the stack with [1][1] then applies fib 10 times).

Additional built in forth functions can be found in the dictionary in `core/words/predefined_words.go`, the predefined words might be a bit confusing since the body of the words are all backwards from the order they would be entered in the repl.

A user can define cusom words with the following syntax:

`: <word_label> a set of space delimited words ;`

for example cube could be:

`: cube dup dup * * ;`

an example of recursion

`: drop2 drop drop ;`
`: dec 1 flip - ;`
`: fibto 0 < IF drop2 dec rotate fib print rotate rotate fibto ELSE drop2 THEN ;`

to use:

`1 1 10` -> will print fibonacci numbers for the specified number of itterations (in this case 10)