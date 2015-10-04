package gmqconf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Constants representing the queue methodology
const (
	USE_MEMORY     = 1
	USE_DATABASE   = 2
	USE_FILESYSTEM = 3
)

type AuthConf struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

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

type LogConf struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ClusterConf struct {
	Proto       string `json:"proto"`
	Port        string `json:"cluster_port"`
	TimeoutMsec int64  `json:"timeout_msec"`
	Cidr        string `json:"cidr"`
}

type Params struct {
	Auth    AuthConf    `json:"auth"`
	Network NetConf     `json:"network"`
	Cluster ClusterConf `json:"cluster"`
	Queue   QConf       `json:"queue"`
	Db      DbConf      `json:"database"`
	Fs      FsConf      `json:"filesystem"`
	Log     LogConf     `json:"log"`
}

// Parse a JSON configuration file
// Params: the string representing the file path
// Return: the pointer to the Param configuration type, an error if the parsing fails
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
