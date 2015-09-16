package gmq

import (
	"sync"
)

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

func (q *Queue) Init(capacity int) *Queue {
	if capacity > 0 {
		q.QObj = make(map[int][]byte, capacity)
	} else {
		q.QObj = make(map[int][]byte, DEFAULT_QUEUE_CAP)
	}
	return q
}

func (q Queue) GetLength() (int, error) {
	ret := len(q.QObj)
	return ret, nil
}

func (q Queue) Create(name string) (QueueInterface, error) {
	mq := Queue{QName: name}
	return mq, nil
}

func (q Queue) Push(o []byte) error {
	last := len(q.QObj)
	q.lock.Lock()
	q.QObj[last] = o
	q.lock.Unlock()
	return nil
}

func (q Queue) Pop() ([]byte, error) {
	var obj []byte
	q.lock.Lock()
	obj = q.QObj[0]
	q.sync()
	q.lock.Unlock()
	return obj, nil
}

func (q Queue) sync() {
	for i := 1; i <= len(q.QObj); i++ {
		q.QObj[i-1] = q.QObj[i]
	}
	delete(q.QObj, len(q.QObj)-1)
}

func (q PrioQueue) Push(prio int, o []byte) {
	last := len(q.QObj[prio])
	if last == 0 {
		tQobj := make(map[int][]byte, DEFAULT_QUEUE_CAP)
		q.QObj[prio] = tQobj
	}
	q.lock.Lock()
	q.QObj[prio][last] = o
	q.lock.Unlock()
}

func (q PrioQueue) Pop(prio int) []byte {
	var obj []byte
	q.lock.Lock()
	if _, ok := q.QObj[prio][0]; !ok {
		q.sync(prio)
	}
	obj = q.QObj[prio][0]
	q.lock.Unlock()
	return obj
}

func (q PrioQueue) sync(prio int) {
	for i := 0; i < len(q.QObj[prio])+1; i++ {
		q.QObj[prio][i] = q.QObj[prio][i+1]
	}
	q.QObj[prio][len(q.QObj[prio])] = nil
}
