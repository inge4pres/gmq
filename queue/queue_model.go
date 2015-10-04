package gmq

import (
	"errors"
)

const (
	DEFAULT_QUEUE_CAP  = 4096
	MAX_QUEUE_NUMBER   = 4096
	MAX_MESSAGE_LENGHT = 40960
)

var QueueInstance map[string]QueueInterface

//QueueInterface: the interface type describing the operations a queue must have
type QueueInterface interface {
	GetLength() (int, error)
	Create(string) (QueueInterface, error)
	Push([]byte) error
	Pop() ([]byte, error)
	sync()
}

//InitQueueInstance: init the singleton map used by the server
func InitQueueInstance(dim int) map[string]QueueInterface {
	QueueInstance = make(map[string]QueueInterface, dim)
	return QueueInstance
}

//GetQueue: return the corresponding named queue in the QueueInstance map
// return nil if no queue with given name is present
func GetQueue(name string) QueueInterface {
	if ret, ok := QueueInstance[name]; ok {
		return ret
	}
	return nil
}

//GetLenght: return the lenght of the queue if it exists in queueInstance or -1, err
func GetQueueLength(qname string) (int, error) {
	queue := GetQueue(qname)
	if queue != nil {
		return queue.GetLength()
	}
	return -1, errors.New("No such queue named " + qname)
}
