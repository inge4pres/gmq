package gmqconf

import (
	"encoding/json"
	"errors"
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

type DbConf struct {
	DbVendor string `json:"vendor"`
	Dsn      string `json:"dsn"`
}

type FsConf struct {
	Path string `json:"path"`
}

type Params struct {
	Network NetConf `json:"network"`
	Queue   QConf   `json:"queue"`
	Db      DbConf  `json:"database"`
	Fs      FsConf  `json:"filesystem"`
}

func ParseConfiguration(file string) (*Params, error) {
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
