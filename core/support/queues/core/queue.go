package core

type DequeueHandler interface {
	HasValue(i int64)
	Empty()
}
type Queue interface {
	Enqueue(i int64)
	Dequeue(handler DequeueHandler)
	Walk(func(i int64))
}

type QueueNode interface {
	GetValue(handler DequeueHandler)

	Previous() QueueNode
	Next() QueueNode

	UpdateNext(QueueNode)
	UpdatePrevious(QueueNode)
}

type StartNode interface {
	QueueNode
}
type EndNode interface {
	QueueNode
}

type queue struct {
	add         StartNode
	remove      EndNode
	nodeFactory ValueNodeFactory
}

func (q *queue) Enqueue(i int64) {
	followingNode := q.add.Next()
	headNode := q.add

	newNodeConnections := NodeConnections{
		Next:     followingNode,
		Previous: q.add,
	}
	newNode := q.nodeFactory.NewInstance(i, newNodeConnections)

	headNode.UpdateNext(newNode)
	followingNode.UpdatePrevious(newNode)
}
func (q *queue) Dequeue(handler DequeueHandler) {
	valueNode := q.remove.Previous()

	tailNode := q.remove
	beforeValueNode := valueNode.Previous()

	tailNode.UpdatePrevious(beforeValueNode)
	beforeValueNode.UpdateNext(tailNode)

	valueNode.GetValue(handler)
}

type walkDequeueHandler struct {
	ifFound    func(i int64)
	ifNotFound func()
}

func (h walkDequeueHandler) Empty() {
	h.ifNotFound()
}
func (h walkDequeueHandler) HasValue(i int64) {
	h.ifFound(i)
}
func (q *queue) Walk(visitor func(i int64)) {
	node := q.remove.Previous()
	finished := false
	for {
		node.GetValue(walkDequeueHandler{
			ifFound:    visitor,
			ifNotFound: func() { finished = true },
		})

		node = node.Previous()

		if finished {
			break
		}
	}
}

type StartNodeFactory interface {
	NewInstance() StartNode
}
type EndNodeFactory interface {
	NewInstance() EndNode
}

type NodeConnections struct {
	Next     QueueNode
	Previous QueueNode
}
type ValueNodeFactory interface {
	NewInstance(value int64, connections NodeConnections) QueueNode
}

func NewQueue(
	startNodeFactory StartNodeFactory,
	endNodeFactory EndNodeFactory,
	valueNodeFactory ValueNodeFactory,
) Queue {
	start := startNodeFactory.NewInstance()
	end := endNodeFactory.NewInstance()

	start.UpdateNext(end)
	end.UpdatePrevious(start)

	return &queue{
		remove:      end,
		add:         start,
		nodeFactory: valueNodeFactory,
	}
}
