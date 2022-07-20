package contextgraph

import (
	"container/list"
	"sync"
)

type queue struct {
	list *list.List
	mx   sync.Mutex
}

func NewJobQueue() JobQueue {
	return &queue{
		list: list.New(),
		mx:   sync.Mutex{},
	}
}

func (q *queue) Push(job ImportJob) {
	q.mx.Lock()
	defer q.mx.Unlock()

	q.list.PushBack(job)
}

func (q *queue) Pop() ImportJob {
	q.mx.Lock()
	defer q.mx.Unlock()

	e := q.list.Front()
	if e != nil {
		q.list.Remove(e)

		return e.Value.(ImportJob)
	}

	return ImportJob{}
}
