package gmq

import (
	"sync"
	"time"
)

const (
	DEFAULT_QUEUE_CAP = 4096
	NO_PRIORITY       = 0
	LOW_PRIORITY      = 1
	MEDIUM_PRIORITY   = 2
	HIGH_PRIORITY     = 3
	MAX_PRIORITY      = 4
)

type QueueInterface interface {
	Push([]byte) error
	Pop() ([]byte, error)
	sync()
}

type QueueManager struct {
	Obj  map[string]*Queue
	Tick time.Timer
}

type Queue struct {
	lock  sync.RWMutex
	QName string
	QObj  map[int][]byte
}

type PrioQueue struct {
	lock  sync.RWMutex
	QName string
	QObj  map[int]map[int][]byte
}

func (q *Queue) Push(o []byte) error {
	last := len(q.QObj)
	if last == 0 {
		q.QObj = make(map[int][]byte, DEFAULT_QUEUE_CAP)
	}
	q.lock.Lock()
	q.QObj[last] = o
	q.lock.Unlock()
	return nil
}

func (q *Queue) Pop() ([]byte, error) {
	var obj []byte
	q.lock.Lock()
	if _, ok := q.QObj[0]; !ok {
		q.sync()
	}
	obj = q.QObj[0]
	q.sync()
	q.lock.Unlock()
	return obj, nil
}

func (q *Queue) sync() {
	for i := 1; i < len(q.QObj); i++ {
		q.QObj[i-1] = q.QObj[i]
	}
}

func (q *PrioQueue) Push(prio int, o []byte) {
	last := len(q.QObj[prio])
	if last == 0 {
		tQobj := make(map[int][]byte, DEFAULT_QUEUE_CAP)
		q.QObj[prio] = tQobj
	}
	q.lock.Lock()
	q.QObj[prio][last] = o
	q.lock.Unlock()
}

func (q *PrioQueue) Pop(prio int) []byte {
	var obj []byte
	q.lock.Lock()
	if _, ok := q.QObj[prio][0]; !ok {
		q.sync(prio)
	}
	obj = q.QObj[prio][0]
	q.lock.Unlock()
	return obj
}

func (q *PrioQueue) sync(prio int) {
	for i := 1; i < len(q.QObj[prio]); i++ {
		q.QObj[prio][i-1] = q.QObj[prio][i]
	}
}
