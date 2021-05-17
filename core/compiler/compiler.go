package compiler

import (
	"fmt"
	"strings"
	"tim/forth/core/support/queues"
	"tim/forth/core/support/queues/core"
	"tim/forth/core/support/stacks"
)

type AddResult int

type ForthCompiler interface {
	PushWord(word string) error
	Complete() (map[string]func(*stacks.ForthStack, stacks.StringStack) error, error)
}

type ExpressionPushHandler interface {
	onRejected()
	onComplete() // Probably should pass in some data here
	onError(error)
}
type CompletionHandler interface {
	onNativeComplete(string, func(*stacks.ForthStack, stacks.StringStack) error)
	onError(string)
}
type ExpressionAccumulator interface {
	push(word string, resultHandler ExpressionPushHandler)
	attemptComplete(completionHandler CompletionHandler)
	toString() string
	name() string
	id() string
}

type ExpressionTokenConsumer interface {
	consume(word string, acc *baseAccumulator) ExpressionTokenConsumer
}

type DequeueToSliceHandler struct {
	isEmpty bool
	index   int64
	result  []string
}

func (h *DequeueToSliceHandler) HasValue(s string) {
	h.result[h.index] = s

	h.index = h.index + 1
}
func (h *DequeueToSliceHandler) Empty() {
	h.isEmpty = true
}
func toSlice(queue ExpressionQueue) []string {
	handler := &DequeueToSliceHandler{
		index:   0,
		isEmpty: false,
		result:  make([]string, queue.getCount()),
	}
	for {
		if handler.isEmpty {
			break
		}
		queue.Dequeue(handler)
	}

	return handler.result
}
func reverse(strings []string) []string {
	result := make([]string, len(strings))
	max := len(strings) - 1
	for index, entry := range strings {
		result[max-index] = entry
	}

	return result
}

type baseAccumulator struct {
	hasLabel   bool
	label      string
	identifier string
	body       ExpressionQueue
	consumer   ExpressionTokenConsumer
}

func (acc *baseAccumulator) push(word string, resultHandler ExpressionPushHandler) {
	acc.consumer = acc.consumer.consume(word, acc)
}
func (acc *baseAccumulator) attemptComplete(handler CompletionHandler) {
	if acc.hasLabel {
		native := wrapPredefined(acc.name(), reverse(toSlice(acc.body)))
		handler.onNativeComplete(acc.name(), native)
		return
	}

	handler.onError("Not complete")
}
func (acc *baseAccumulator) toString() string {
	return fmt.Sprintf("[%s] -> %s", acc.label, acc.body.ToString())
}

func (acc *baseAccumulator) name() string {
	return acc.label
}

func (acc *baseAccumulator) id() string {
	return acc.identifier
}

type bodyConsumer struct{}

func (b bodyConsumer) consume(word string, acc *baseAccumulator) ExpressionTokenConsumer {
	acc.body.Enqueue(word)
	return b
}

type labelConsumer struct{}

func (b labelConsumer) consume(word string, acc *baseAccumulator) ExpressionTokenConsumer {
	acc.label = word
	acc.hasLabel = true

	return &bodyConsumer{}
}

func NewBaseAccumulator(id string) ExpressionAccumulator {
	return &baseAccumulator{
		identifier: id,
		consumer:   &labelConsumer{},
		hasLabel:   false,
		label:      "",
		body:       NewExpressionQueue(),
	}
}

type ifExpressionConsumer interface {
	consume(token string, accume *ifAccumulator, onReject func(), onComplete func())
}

type consumerUtils struct{}

func (util consumerUtils) isThen(s string) bool {
	return strings.ToLower(s) == "then"
}
func (util consumerUtils) isElse(s string) bool {
	return strings.ToLower(s) == "else"
}

type elseConsumer struct {
	utils consumerUtils
}

func (consumer elseConsumer) consume(token string, accume *ifAccumulator, onReject func(), onComplete func()) {
	if consumer.utils.isThen(token) {
		onComplete()
		return
	}

	if consumer.utils.isElse(token) {
		onReject()
		return
	}

	accume.elseBody.Enqueue(token)
}

type ifConsumer struct {
	utils consumerUtils
}

func (consumer ifConsumer) consume(token string, accume *ifAccumulator, onReject func(), onComplete func()) {
	if consumer.utils.isThen(token) {
		onComplete()
		return
	}

	if consumer.utils.isElse(token) {
		accume.consumer = &elseConsumer{
			utils: consumer.utils,
		}
		return
	}

	accume.ifBody.Enqueue(token)
}

type ifAccumulator struct {
	label      string
	identifier string
	ifBody     ExpressionQueue
	elseBody   ExpressionQueue

	isComplete bool

	consumer ifExpressionConsumer
}

func (acc *ifAccumulator) push(word string, resultHandler ExpressionPushHandler) {
	acc.consumer.consume(word, acc, resultHandler.onRejected, func() {
		acc.isComplete = true
		resultHandler.onComplete()
	})
}

func (acc *ifAccumulator) attemptComplete(handler CompletionHandler) {
	if acc.isComplete {
		ifList := reverse(toSlice(acc.ifBody))
		elseList := reverse(toSlice(acc.elseBody))

		handler.onNativeComplete(acc.label, func(forthStack *stacks.ForthStack, executionStack stacks.StringStack) error {
			if forthStack.IsEmpty() {
				fmt.Println("Underflow....")
				return nil
			}

			if forthStack.Peek().ValueOf() == 0 {
				for _, v := range ifList {
					executionStack.Push(v)
				}
			} else {
				for _, v := range elseList {
					executionStack.Push(v)
				}
			}

			return nil
		})
	}
	handler.onError("Can not complete an if handler without a then")
}
func (acc *ifAccumulator) toString() string {
	return fmt.Sprintf("[%s] -> [IF] %s [ELSE] %s \nisComplete -> %t\n", acc.name(), acc.ifBody.ToString(), acc.elseBody.ToString(), acc.isComplete)
}
func (acc *ifAccumulator) name() string {
	return acc.label
}
func (acc *ifAccumulator) id() string {
	return acc.identifier
}

func NewIfExpressionAccumulator(id string) ExpressionAccumulator {
	return &ifAccumulator{
		identifier: id,
		label:      id,
		ifBody:     NewExpressionQueue(),
		elseBody:   NewExpressionQueue(),
		consumer:   &ifConsumer{utils: consumerUtils{}},
	}
}

type ExpressionBuilder interface {
	add(word string) AddResult
	name() string
	expBody() []string
	isNull() bool
	toString() string
}

type ExpressionDequeueHandler interface {
	HasValue(s string)
	Empty()
}
type ExpressionQueue interface {
	Enqueue(s string)
	Dequeue(handler ExpressionDequeueHandler)
	getCount() int64
	ToString() string
}

type expressionQueue struct {
	backingQueue core.Queue
	keyValueMap  map[int64]string
	count        int64
}

func (q *expressionQueue) Enqueue(s string) {
	count := q.count

	q.backingQueue.Enqueue(int64(count))
	q.keyValueMap[count] = s

	q.count = count + 1
}

type underlyingQueueDequeueHandler struct {
	q       *expressionQueue
	handler ExpressionDequeueHandler
}

func (u *underlyingQueueDequeueHandler) Empty() {
	u.handler.Empty()
}
func (u *underlyingQueueDequeueHandler) HasValue(i int64) {
	result := u.q.keyValueMap[i]

	u.handler.HasValue(result)

	delete(u.q.keyValueMap, i)
	u.q.count = u.q.count - 1
}

func (q *expressionQueue) Dequeue(handler ExpressionDequeueHandler) {
	q.backingQueue.Dequeue(&underlyingQueueDequeueHandler{q: q, handler: handler})
}
func (q *expressionQueue) getCount() int64 {
	return q.count
}
func (q *expressionQueue) ToString() string {
	entries := make([]int64, q.count)
	index := 0
	q.backingQueue.Walk(func(i int64) {
		entries[index] = i
		index = index + 1
	})

	result := ""
	for _, v := range entries {
		value := q.keyValueMap[v]
		result = fmt.Sprintf("%s[%s]", result, value)
	}

	return result
}

func NewExpressionQueue() ExpressionQueue {
	queue := queues.NewInstance()

	return &expressionQueue{
		backingQueue: queue,
		count:        0,
		keyValueMap:  make(map[int64]string),
	}
}

type pushHandler struct {
	c *forthCompiler
}

func (h *pushHandler) onComplete() {
	c := h.c
	referenceWord := c.currentExpression.name()

	poppedExpressionId := c.expressionIdStack.Pop()
	poppedExpression := c.expressionMap[poppedExpressionId]
	fmt.Println(poppedExpression.toString())

	c.currentExpression = poppedExpression
	c.currentExpression.push(referenceWord, h)
}

func (h *pushHandler) onError(e error) {

}
func (h *pushHandler) onRejected() {}

func NewResultHandler(c *forthCompiler) ExpressionPushHandler {
	return &pushHandler{
		c: c,
	}
}

func wrapPredefined(name string, body []string) func(*stacks.ForthStack, stacks.StringStack) error {
	return func(forthStack *stacks.ForthStack, executionStack stacks.StringStack) error {
		for _, w := range body {
			executionStack.Push(w)
		}
		return nil
	}
}

type completionResult struct {
	hasError        bool
	errorMessage    string
	nativeFunctions map[string]func(*stacks.ForthStack, stacks.StringStack) error
}

func NewCompletionResult() *completionResult {
	return &completionResult{
		hasError:        false,
		nativeFunctions: make(map[string]func(*stacks.ForthStack, stacks.StringStack) error),
	}
}

type completionHandler struct {
	c *completionResult
}

func (h *completionHandler) onError(message string) {

}
func (h *completionHandler) onNativeComplete(label string, nativeFunc func(*stacks.ForthStack, stacks.StringStack) error) {
	h.c.nativeFunctions[label] = nativeFunc
}

func NewCompletionHandler(c *completionResult) CompletionHandler {
	return &completionHandler{
		c: c,
	}
}

type forthCompiler struct {
	idGenerator       InternalIdProvider
	currentExpression ExpressionAccumulator
	expressionIdStack stacks.StringStack
	expressionMap     map[string]ExpressionAccumulator
}

func (c *forthCompiler) PushWord(word string) error {
	if strings.ToLower(word) == "if" {
		exp := NewIfExpressionAccumulator(c.idGenerator.NextId())

		c.expressionIdStack.Push(c.currentExpression.id())

		c.expressionMap[exp.id()] = exp
		c.currentExpression = exp

		return nil
	}

	c.currentExpression.push(word, NewResultHandler(c))
	return nil
}

func (c *forthCompiler) Complete() (map[string]func(*stacks.ForthStack, stacks.StringStack) error, error) {
	result := NewCompletionResult()
	handler := NewCompletionHandler(result)
	for _, expression := range c.expressionMap {
		fmt.Printf("Attempting to complete -> \n%s\n", expression.toString())
		expression.attemptComplete(handler)

		if result.hasError {
			break
		}
	}

	return result.nativeFunctions, nil
}

type InternalIdProvider interface {
	NextId() string
}

func NewCompiler(idGenerator InternalIdProvider) ForthCompiler {
	baseId := idGenerator.NextId()

	baseExpression := NewBaseAccumulator(baseId)
	expressionMap := make(map[string]ExpressionAccumulator)
	expressionIdStack := stacks.NewStringStack()

	expressionMap[baseId] = baseExpression

	return &forthCompiler{
		idGenerator:       idGenerator,
		currentExpression: baseExpression,
		expressionIdStack: expressionIdStack,
		expressionMap:     expressionMap,
	}
}
