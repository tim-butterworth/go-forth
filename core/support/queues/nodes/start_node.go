package nodes

import "tim/forth/core/support/queues/core"

type startNode struct {
	nextNode core.QueueNode
}

func (e startNode) GetValue(handler core.DequeueHandler) {
	handler.Empty()
}
func (e *startNode) Previous() core.QueueNode {
	return e
}
func (e *startNode) Next() core.QueueNode {
	return e.nextNode
}
func (e startNode) UpdatePrevious(previousNode core.QueueNode) {}
func (e *startNode) UpdateNext(nextNode core.QueueNode) {
	e.nextNode = nextNode
}

func NewStartNode() core.QueueNode {
	return &startNode{}
}

type startNodeFactory struct{}

func (factory startNodeFactory) NewInstance() core.StartNode {
	return &startNode{}
}
func StartNodeFactory() core.StartNodeFactory {
	return startNodeFactory{}
}
