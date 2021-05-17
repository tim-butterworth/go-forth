package queues

import (
	"tim/forth/core/support/queues/core"
	"tim/forth/core/support/queues/nodes"
)

func NewInstance() core.Queue {
	return core.NewQueue(
		nodes.StartNodeFactory(),
		nodes.EndNodeFacotry(),
		nodes.ValueNodeFactory(),
	)
}
