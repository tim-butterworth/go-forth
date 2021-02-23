package stacks

type StringNode interface {
	isEmpty() bool
	value() string
	next() StringNode
}
type emptyStringNode struct{}

func (e emptyStringNode) isEmpty() bool {
	return true
}
func (e emptyStringNode) value() string {
	return ""
}
func (e emptyStringNode) next() StringNode {
	return e
}

type stringNode struct {
	v string
	n StringNode
}

func (p stringNode) isEmpty() bool {
	return false
}
func (p stringNode) value() string {
	return p.v
}
func (p stringNode) next() StringNode {
	return p.n
}

type StringStack interface {
	IsEmpty() bool
	Push(string)
	Pop() string
	Size() int
}
type Stack struct {
	root  StringNode
	count int
}

func (stack *Stack) IsEmpty() bool {
	return stack.root.isEmpty()
}
func (stack *Stack) Push(s string) {
	newNode := stringNode{v: s, n: stack.root}
	stack.root = newNode

	stack.count = stack.count + 1
}
func (stack *Stack) Pop() string {
	value := stack.root.value()
	stack.root = stack.root.next()

	if stack.count > 0 {
		stack.count = stack.count - 1
	}
	return value
}
func (stack *Stack) Size() int {
	return stack.count
}

func NewStringStack() StringStack {
	return &Stack{
		root:  emptyStringNode{},
		count: 0,
	}
}
