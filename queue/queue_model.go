package gmq

const (
	DEFAULT_QUEUE_CAP  = 4096
	MAX_QUEUE_NUMBER   = 4096
	MAX_MESSAGE_LENGHT = 40960
)

var QueueInstance map[string]QueueInterface

type QueueInterface interface {
	Create(string) QueueInterface
	Push([]byte) error
	Pop() ([]byte, error)
	sync()
}

func init() {
	QueueInstance = make(map[string]QueueInterface, 1)
}

func GetQueue(name string) (*QueueInterface, bool) {
	ret, ok := QueueInstance[name]
	return &ret, ok
}
