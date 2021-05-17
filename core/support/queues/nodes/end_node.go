package nodes

import "tim/forth/core/support/queues/core"

type endNode struct {
	previousNode core.QueueNode
}

func (e endNode) GetValue(handler core.DequeueHandler) {
	handler.Empty()
}
func (e *endNode) Previous() core.QueueNode {
	return e.previousNode
}
func (e *endNode) Next() core.QueueNode {
	return e
}
func (e *endNode) UpdateNext(nextNode core.QueueNode) {}
func (e *endNode) UpdatePrevious(previousNode core.QueueNode) {
	e.previousNode = previousNode
}

type endNodeFactory struct{}

func (factory endNodeFactory) NewInstance() core.EndNode {
	return &endNode{}
}
func EndNodeFacotry() core.EndNodeFactory {
	return endNodeFactory{}
}
