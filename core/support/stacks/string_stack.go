package stacks

import (
	"fmt"
	"tim/forth/core/support/stacks/internal"
)

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
	return "<-- EMPTY -->"
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
	ToString() string
}
type Stack struct {
	root         StringNode
	count        int
	backingStack internal.BackingStack
	dataMap      map[int64]string
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
func (stack *Stack) ToString() string {
	node := stack.root

	result := ""
	for {
		if node.isEmpty() {
			break
		}

		result = result + fmt.Sprintf("[%s]", node.value())
		node = node.next()
	}

	return result
}

func NewStringStack() StringStack {
	baseStack := internal.NewBackingStack()
	dataMap := make(map[int64]string)

	return &Stack{
		root:         emptyStringNode{},
		count:        0,
		backingStack: baseStack,
		dataMap:      dataMap,
	}
}
