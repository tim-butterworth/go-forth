package stacks

import "fmt"

type ForthItem interface {
	IsEmpty() bool
	ToString() string
	ValueOf() int64
}
type Empty struct{}

func (e Empty) IsEmpty() bool {
	return true
}
func (e Empty) ToString() string {
	return "--empty--"
}
func (e Empty) ValueOf() int64 {
	return 0
}

type Number struct {
	Value int64
}

func (n Number) IsEmpty() bool {
	return false
}
func (n Number) ToString() string {
	return fmt.Sprintf("%d", n.Value)
}
func (n Number) ValueOf() int64 {
	return n.Value
}

type forthNode interface {
	next() forthNode
	isEmpty() bool
	value() ForthItem
}
type emptyNode struct{}

func (e emptyNode) next() forthNode {
	return e
}
func (e emptyNode) isEmpty() bool {
	return true
}
func (e emptyNode) value() ForthItem {
	return Empty{}
}

type populatedNode struct {
	nextNode forthNode
	item     ForthItem
}

func (p *populatedNode) next() forthNode {
	return p.nextNode
}
func (p *populatedNode) isEmpty() bool {
	return false
}
func (p *populatedNode) value() ForthItem {
	return p.item
}

func newNode(next forthNode, item ForthItem) forthNode {
	return &populatedNode{
		nextNode: next,
		item:     item,
	}
}

type ForthStack struct {
	root forthNode
}

func (stack *ForthStack) Push(item ForthItem) {
	node := newNode(stack.root, item)
	stack.root = node
}
func (stack *ForthStack) Pop() ForthItem {
	result := stack.root

	if result == nil {
		return Empty{}
	}
	stack.root = stack.root.next()
	return result.value()
}
func (stack *ForthStack) Peek() ForthItem {
	return stack.root.value()
}
func (stack *ForthStack) IsEmpty() bool {
	return stack.root.isEmpty()
}
func (stack *ForthStack) ToString() string {
	node := stack.root

	result := ""
	for {
		if node.isEmpty() {
			break
		}

		result = result + fmt.Sprintf("[%s]", node.value().ToString())
		node = node.next()
	}

	return result
}

func NewStack() *ForthStack {
	return &ForthStack{root: emptyNode{}}
}
