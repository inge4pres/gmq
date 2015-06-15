package gmq

import (
	"log"
	"os"
	"sync"
)

const (
	DEFAULT_QUEUE_CAP = 1024
	LOW_PRIORITY      = 1
	MEDIUM_PRIORITY   = 2
	HIGH_PRIORITY     = 3
)

var l *log.Logger

type Enqueuer interface {
	Push()
	Sync()
}

type Dequeuer interface {
	Pop()
	Sync()
}

type Queue struct {
	lock   sync.RWMutex
	QName  string
	QTopic []string
	QObj   map[int]interface{}
}

type PrioQueue struct {
	lock   sync.RWMutex
	QName  string
	QTopic []string
	QObj   map[int]map[int]interface{}
}

func init() {
	l := log.New(os.Stdout, "GMQ: ", log.LstdFlags)
	l.Println("Init pkg GMQ")
}

func (q *Queue) Push(args ...interface{}) {
	last := len(q.QObj)
	if last == 0 {
		q.QObj = make(map[int]interface{})
	}
	for o := range args {
		q.lock.Lock()
		q.QObj[last] = o
		q.lock.Unlock()
	}
	l.Printf("Queue \"%s\" now has %d elements", q.QName, len(q.QObj))
}

func (q *Queue) Pop() interface{} {
	//TODO
	return nil
}
