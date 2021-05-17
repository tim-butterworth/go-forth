package compiler

import (
	"fmt"
	"tim/forth/core/support/queues"
	"tim/forth/core/support/queues/core"
	"tim/forth/core/support/stacks"
)

type AddResult int

const (
	SUCCESS AddResult = iota
	ERROR
	EXPRESSION_COMPLETE
)

type ForthCompiler interface {
	PushWord(word string) error
	Complete() (error, map[string][]string)
}

type ExpressionBuilder interface {
	add(word string) AddResult
	name() string
	expBody() []string
	isNull() bool
	toString() string
}

type baseExpression struct {
	label string
	body  ExpressionQueue
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

func (exp *baseExpression) add(word string) AddResult {
	exp.body.Enqueue(word)

	return SUCCESS
}
func (exp *baseExpression) name() string {
	return exp.label
}

type dequeueHandler struct {
	ifFound func(s string)
}

func (h dequeueHandler) Empty() {}
func (h dequeueHandler) HasValue(s string) {
	h.ifFound(s)
}
func (exp *baseExpression) expBody() []string {
	bodyStack := exp.body
	result := make([]string, bodyStack.getCount())
	for i := range result {
		bodyStack.Dequeue(&dequeueHandler{
			ifFound: func(s string) {
				result[i] = s
			},
		})
	}

	return result
}
func (exp *baseExpression) isNull() bool {
	return false
}
func (exp *baseExpression) toString() string {
	return fmt.Sprintf("BaseExpression [%s] -> %s", exp.label, exp.body.ToString())
}
func NewBaseExpression(name string) *baseExpression {
	return &baseExpression{
		label: name,
		body:  NewExpressionQueue(),
	}
}

type ifExpression struct {
	id       string
	isIf     bool
	ifBody   ExpressionQueue
	elseBody ExpressionQueue
}

func (exp *ifExpression) add(word string) AddResult {
	if word == "then" {
		return EXPRESSION_COMPLETE
	}

	if exp.isIf {
		if word == "else" {
			exp.isIf = false
		} else {
			exp.ifBody.Enqueue(word)
		}
	} else {
		if word == "else" {
			return ERROR
		}

		exp.elseBody.Enqueue(word)
	}

	return SUCCESS
}
func (exp *ifExpression) name() string {
	return exp.id
}
func (exp *ifExpression) expBody() []string {
	return make([]string, 0)
}
func (exp *ifExpression) isNull() bool {
	return false
}
func (exp *ifExpression) toString() string {
	ifBody := exp.ifBody
	elseBody := exp.elseBody

	return fmt.Sprintf(
		"IF_EXPRESSION [%s] -> \n\tIF: %s\n\tELSE: %s",
		exp.id,
		ifBody.ToString(),
		elseBody.ToString(),
	)
}

func NewIfExpression(id string) *ifExpression {
	return &ifExpression{
		id:       id,
		isIf:     true,
		ifBody:   NewExpressionQueue(),
		elseBody: NewExpressionQueue(),
	}
}

type nullExpression struct{}

func (exp nullExpression) add(word string) AddResult {
	return ERROR
}
func (exp nullExpression) name() string {
	return ""
}
func (exp nullExpression) isNull() bool {
	return true
}
func (exp nullExpression) toString() string {
	return "<--- NULL_EXPRESSION --->"
}
func (exp nullExpression) expBody() []string {
	return make([]string, 0)
}

type forthCompiler struct {
	idGenerator       InternalIdProvider
	currentExpression ExpressionBuilder
	expressionIdStack stacks.StringStack
	expressionMap     map[string]ExpressionBuilder
}

func (c *forthCompiler) PushWord(word string) error {
	if c.currentExpression.isNull() {
		baseExpression := NewBaseExpression(word)
		c.expressionMap[word] = baseExpression
		c.currentExpression = baseExpression
	}

	if word == "if" {
		exp := NewIfExpression(c.idGenerator.NextId())

		c.expressionIdStack.Push(c.currentExpression.name())

		c.expressionMap[exp.name()] = exp
		c.currentExpression = exp

		return nil
	}

	result := c.currentExpression.add(word)

	if result == EXPRESSION_COMPLETE {
		referenceWord := c.currentExpression.name()

		poppedExpressionId := c.expressionIdStack.Pop()
		poppedExpression := c.expressionMap[poppedExpressionId]

		fmt.Println(c.currentExpression.toString())
		c.currentExpression = poppedExpression
		c.currentExpression.add(referenceWord)

		return nil
	}

	if result == ERROR {
		fmt.Println("ERROROROROROROR!")
		return nil //this should be an error
	}

	return nil
}
func (c *forthCompiler) Complete() (error, map[string][]string) {
	fmt.Println(c.currentExpression.toString())
	result := make(map[string][]string)
	for _, expression := range c.expressionMap {
		result[expression.name()] = expression.expBody()
	}
	return nil, result
}

type InternalIdProvider interface {
	NextId() string
}

func NewCompiler(idGenerator InternalIdProvider) ForthCompiler {
	return &forthCompiler{
		idGenerator:       idGenerator,
		currentExpression: nullExpression{},
		expressionIdStack: stacks.NewStringStack(),
		expressionMap:     make(map[string]ExpressionBuilder),
	}
}
