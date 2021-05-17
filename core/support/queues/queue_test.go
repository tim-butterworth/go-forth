package queues_test

import (
	"fmt"
	"testing"
	"tim/forth/core/support/queues"
	"tim/forth/core/support/queues/core"
)

func Test_queueStartsEmpty(t *testing.T) {
	queue := queues.NewInstance()

	testHandler := NewTestDequeueHandler(func(i int64) {}, func() {})

	queue.Dequeue(&testHandler)

	if testHandler.hasValueCalled {
		t.Error("[HasValue] should not have been called")
		return
	}

	if !testHandler.emptyCalled {
		t.Error("[Empty] should have been called")
	}
}

func Test_canEnqueue_and_Dequeue_a_value(t *testing.T) {
	queue := queues.NewInstance()

	expected := int64(42)
	queue.Enqueue(expected)

	testHandler := NewTestDequeueHandler(
		func(actual int64) {
			if actual != expected {
				t.Errorf("Expected %d but got %d", expected, actual)
			}
		},
		func() {},
	)

	queue.Dequeue(&testHandler)

	if testHandler.emptyCalled {
		t.Error("Empty should not have been called")
		return
	}

	if !testHandler.hasValueCalled {
		t.Error("HasValue should have been called")
	}
}

func Test_isEmpty_afterAllDequed(t *testing.T) {
	queue := queues.NewInstance()

	queue.Enqueue(int64(12))

	valueHandler := NewTestDequeueHandler(func(i int64) {}, func() {})
	queue.Dequeue(&valueHandler)

	emptyHandler := NewTestDequeueHandler(func(i int64) {}, func() {})
	queue.Dequeue(&emptyHandler)

	if emptyHandler.hasValueCalled {
		t.Error("Should not have a value")
		return
	}

	if !emptyHandler.emptyCalled {
		t.Error("Empty should have been called")
	}
}

func Test_canEnqueue_and_Dequeue_multipleValues(t *testing.T) {
	queue := queues.NewInstance()
	values := []int64{
		1,
		2,
		3,
		5,
		7,
		11,
		13,
		17,
	}

	for _, v := range values {
		queue.Enqueue(v)
	}

	var returnedValue int64
	valueHandler := NewTestDequeueHandler(
		func(i int64) {
			returnedValue = i
		},
		func() {
		},
	)
	for _, v := range values {
		queue.Dequeue(&valueHandler)

		if returnedValue != v {
			t.Errorf("Expected %d but got %d", v, returnedValue)
			break
		} else {
			fmt.Printf("Matched %d -> %d\n", v, returnedValue)
		}

		if valueHandler.emptyCalled {
			t.Error("Empty should not be called")
			break
		}
	}

	queue.Dequeue(&valueHandler)

	if !valueHandler.emptyCalled {
		t.Error("The queue should have been empty")
	}
}

func Test_canFill_empty_refill_empty(t *testing.T) {
	queue := queues.NewInstance()

	values1 := []int64{2, 3, 5, 7, 11}
	fillQueue(queue, values1)
	emptyTheQueue(t, queue, values1)

	values2 := []int64{13, 17, 19, 23}
	fillQueue(queue, values2)
	emptyTheQueue(t, queue, values2)
}

func Test_can_fill_partEmpty_fillMore_empty(t *testing.T) {
	queue := queues.NewInstance()

	values1 := []int64{2, 3, 5, 7}
	values2 := []int64{-2, -3, -5, -7}
	values3 := []int64{11, 13, 17, 19, 23, 29, 31}

	fillQueue(queue, values1)
	fillQueue(queue, values2)
	fillQueue(queue, values2)

	removeFromQueue(t, queue, values1)
	removeFromQueue(t, queue, values2)

	fillQueue(queue, values3)

	removeFromQueue(t, queue, values2)
	emptyTheQueue(t, queue, values3)
}

func Test_walk_does_not_invoke_the_callback_for_empty_queue(t *testing.T) {
	queue := queues.NewInstance()

	queue.Walk(func(i int64) {
		t.Error("Walk should not have invoked the callback for an empty queue")
	})
}

func Test_can_retrieve_all_values_with_walk(t *testing.T) {
	queue := queues.NewInstance()
	values := []int64{-1, 11, 23, 31, 43, 51}

	fillQueue(queue, values)

	visited := make([]int64, len(values))
	index := 0
	queue.Walk(func(i int64) {
		visited[index] = i
		index = index + 1
	})

	for i, v := range values {
		if visited[i] != v {
			t.Errorf("Expected: [%d] but got: [%d]", v, visited[i])
		}
	}
}

func Test_walk_does_not_modify_queue(t *testing.T) {
	queue := queues.NewInstance()
	values := []int64{-1, 11, 23, 31, 43, 51}

	fillQueue(queue, values)

	visited := make([]int64, len(values))
	index := 0
	queue.Walk(func(i int64) {
		visited[index] = i
		index = index + 1
	})

	for i, v := range values {
		if visited[i] != v {
			t.Errorf("Expected: [%d] but got: [%d]", v, visited[i])
		}
	}

	emptyTheQueue(t, queue, visited)
}

func fillQueue(queue core.Queue, values []int64) {
	for _, value := range values {
		queue.Enqueue(value)
	}
}

func removeFromQueue(t *testing.T, queue core.Queue, expectedValues []int64) {
	var actual int64
	dequeueHander := NewTestDequeueHandler(
		func(i int64) {
			actual = i
		},
		func() {},
	)

	for _, expectedValue := range expectedValues {
		queue.Dequeue(&dequeueHander)

		if dequeueHander.emptyCalled {
			t.Errorf("The queue should not be empty, expected value %d\n", expectedValue)
			break
		}

		if expectedValue != actual {
			t.Errorf("Expected %d but got %d\n", expectedValue, actual)
			break
		}
	}
}
func emptyTheQueue(t *testing.T, queue core.Queue, expectedValues []int64) {
	removeFromQueue(t, queue, expectedValues)

	dequeueHander := NewTestDequeueHandler(func(i int64) {}, func() {})

	queue.Dequeue(&dequeueHander)
	if !dequeueHander.emptyCalled {
		t.Error("Queue should be empty now")
	}
}

func NewTestDequeueHandler(
	hasValueFun func(i int64),
	emptyFun func(),
) testDequeueHandler {
	return testDequeueHandler{
		hasValueFun:    hasValueFun,
		hasValueCalled: false,

		emptyFun:    emptyFun,
		emptyCalled: false,
	}
}

type testDequeueHandler struct {
	hasValueFun    func(i int64)
	hasValueCalled bool
	hasValueParam  int64
	emptyFun       func()
	emptyCalled    bool
}

func (h *testDequeueHandler) HasValue(i int64) {
	h.hasValueCalled = true
	h.hasValueFun(i)
}

func (h *testDequeueHandler) Empty() {
	h.emptyCalled = true
	h.emptyFun()
}
