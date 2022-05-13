package main

import (
	"sync/atomic"
)

type queue struct {
	placed atomic.Value
}

func NewQueue() *queue {
	que := new(queue)
	que.placed.Store(int(0))

	return que
}

func (q *queue) AppendRecord() {
	q.placed.Store(q.placed.Load().(int) + 1)
}

func (q *queue) DepartRecord() {
	q.placed.Store(q.placed.Load().(int) - 1)
}

func (q *queue) Length() int {
	return q.placed.Load().(int)
}
