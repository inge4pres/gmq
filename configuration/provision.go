package gmqconf

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/inge4pres/gmq/queue"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
	DEFAULT_PROTOCOL    = "tcp"
	DEFAULT_INET        = ""
)

var server *Server
var cluster map[string]*Server

type Server struct {
	Proto, LocalInet, Port string
	Listener               net.Listener
}

func init() {
	server = new(Server)
	cluster = make(map[string]*Server)
}

// Configure the queue used in the server, handling mandatory configurations
// Params: the pointer to the configuration Params
// Returns: the QueueInterface type that will be used in server, an error if configuration fails
func ConfigureQueue(conf *Params) (queue gmq.QueueInterface, err error) {

	if conf.Queue.MaxQueueN < 1 {
		err = errors.New("Please configure MAX_QUEUE_NUMBER with a positive number")
	}
	if conf.Queue.MaxQueueC < 1 {
		err = errors.New("Please configure MAX_QUEUE_CAPACITY with a positive number")
	}
	if conf.Queue.MaxMessageL < 1 {
		err = errors.New("Please configure MAX_MESSAGE_LENGHT with a positive number")
	}

	switch conf.Queue.QueueType {
	case USE_MEMORY:
		mq := gmq.Queue{}
		mq.Init(conf.Queue.MaxQueueC)
		queue = mq

	case USE_DATABASE:
		queue = gmq.DbQueue{}

	case USE_FILESYSTEM:
		queue = gmq.FsQueue{}

	default:
		err = errors.New("Please configure QUEUE_TYPE with 1 (memory), 2 (database) or 3 (filesystem)")
	}

	return queue, err
}

func configureServer(conf *Params) *Server {
	var inet, proto, port string
	if conf.Network.Port == "" {
		port = DEFAULT_LISTEN_PORT
	} else {
		port = conf.Network.Port
	}
	if conf.Network.Proto == "" {
		proto = DEFAULT_PROTOCOL
	} else {
		proto = conf.Network.Proto
	}
	if conf.Network.Inet == "" {
		inet = DEFAULT_INET
	} else {
		inet = conf.Network.Inet
	}
	return &Server{
		Port:      port,
		Proto:     proto,
		LocalInet: inet,
	}
}

// Init the gmqserver singleton with configuration options
// Params: the pointer to a configuration Params type
// Returns: the pointer to the Server, an error if init fails
func InitServer(params *Params) (server *Server, err error) {
	server = configureServer(params)
	server.Listener, err = net.Listen(server.Proto, server.LocalInet+":"+server.Port)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// Stops the running TCP server
func (server *Server) StopServer() {
	server.Listener.Close()
}

// Get the JSON byte array decribing the peers in cluster
func GetPeerList() ([]byte, error) {
	return json.Marshal(cluster)
}
