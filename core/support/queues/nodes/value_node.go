package nodes

import "tim/forth/core/support/queues/core"

type valueNode struct {
	value        int64
	previousNode core.QueueNode
	nextNode     core.QueueNode
}

func (v *valueNode) GetValue(handler core.DequeueHandler) {
	handler.HasValue(v.value)
}
func (v valueNode) IsEmpty() bool {
	return false
}
func (v *valueNode) Previous() core.QueueNode {
	return v.previousNode
}
func (v *valueNode) Next() core.QueueNode {
	return v.nextNode
}
func (v *valueNode) UpdateNext(nextNode core.QueueNode) {
	v.nextNode = nextNode
}
func (v *valueNode) UpdatePrevious(previousNode core.QueueNode) {
	v.previousNode = previousNode
}

type valueNodeFactory struct{}

func (factory valueNodeFactory) NewInstance(value int64, connections core.NodeConnections) core.QueueNode {
	return &valueNode{
		value:        value,
		previousNode: connections.Previous,
		nextNode:     connections.Next,
	}
}
func ValueNodeFactory() core.ValueNodeFactory {
	return valueNodeFactory{}
}
