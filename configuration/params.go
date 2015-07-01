package gmqconf

import (
	"encoding/json"
	"errors"
	comm "gmq/communication"
	q "gmq/queue"
	"io/ioutil"
)

const (
	USE_MEMORY     = 1
	USE_DATABASE   = 2
	USE_FILESYSTEM = 3
)

type NetConf struct {
	Port  string `json:"port"`
	Inet  string `json:"interface"`
	Proto string `json:"protocol"`
}

type QConf struct {
	QueueType   int `json:"queue_type"`
	MaxQueueN   int `json:"max_queue_number"`
	MaxMessageL int `json:"max_message_length"`
	MaxQueueC   int `json:"max_queue_capacity"`
}

type Params struct {
	Network NetConf `json:"network"`
	Queue   QConf   `json:"queue"`
}

func parseConfiguration(file string) (*Params, error) {
	params := new(Params)
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.New("Configuration file not available")
	}
	if err := json.Unmarshal(conf, params); err != nil {
		return nil, errors.New("Unable to parse configuration file, check JSON is well formed")
	}
	return params, err
}

func ConfigureServer(conf *Params) *comm.Server {
	s := new(comm.Server)
	if conf.Network.Port == "" {
		s.Port = comm.DEFAULT_LISTEN_PORT
	}
	if conf.Network.Proto == "" {
		s.Proto = comm.DEFAULT_PROTOCOL
	}
	return s
}

func ConfigureQueue(conf *Params) (*q.QueueInterface, error) {
	queue := new(q.QueueInterface)
	if conf.Queue.MaxQueueN < 1 {
		return nil, errors.New("Please configure MAX_QUEUE_NUMBER with a positive number")
	}
	if conf.Queue.MaxQueueC < 1 {
		return nil, errors.New("Please configure MAX_QUEUE_CAPACITY with a positive number")
	}
	if conf.Queue.MaxMessageL < 1 {
		return nil, errors.New("Please configure MAX_MESSAGE_LENGHT with a positive number")
	}
	switch conf.Queue.QueueType {
	case USE_MEMORY:
		queue := new(q.Queue)
	case USE_DATABASE:
		queue := new(q.DbQueue)
	case USE_FILESYSTEM:
		queue := new(q.FsQueue)
	default:
		return nil, errors.New("Please configure QUEUE_TYPE")
	}
	return queue, nil
}
