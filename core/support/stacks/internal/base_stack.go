package internal

type StackNode interface {
	isEmpty() bool
	value(hasValue func(int64), hasNoValue func())
	next() StackNode
}

type emptyNode struct{}

func (e emptyNode) isEmpty() bool {
	return true
}
func (e emptyNode) value(hasValue func(int64), hasNoValue func()) {
	hasNoValue()
}
func (e emptyNode) next() StackNode {
	return e
}

type node struct{
	nextNode StackNode
	data int64
}

func (n *node) isEmpty() bool {
	return false
}
func (n *node) value(hasValue func(int64), hasNoValue func()) {
	hasValue(n.data)
}
func (n *node) next() StackNode {
	return n.nextNode
}

func newEmptyNode() StackNode {
	return &emptyNode{}
}
func newNode(nextNode StackNode, value int64) StackNode {
	return &node{
		nextNode: nextNode,
		data: value,
	}
}

type BackingStack interface {
	IsEmpty() bool
	Push(i int64)
	Pop(hasValue func(int64), hasNoValue func())
}

type backingStack struct {
	root StackNode
}

func (stack *backingStack) IsEmpty() bool {
	return stack.root.isEmpty()
}
func (stack *backingStack) Push(i int64) {
	newNode := newNode(stack.root, i)

	stack.root = newNode
}
func (stack *backingStack) Pop(hasValue func(int64), hasNoValue func()) {
	node := stack.root

	stack.root = stack.root.next()

	node.value(hasValue, hasNoValue)
}

func NewBackingStack() BackingStack {
	return &backingStack{
		root: newEmptyNode(),
	}
}
