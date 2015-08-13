package gmqnet

import (
	"encoding/base64"
	"errors"
	m "gmq/configuration"
	q "gmq/queue"
	"net"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
	DEFAULT_PROTOCOL    = "tcp"
	DEFAULT_INET        = ""
)

var server *Server

type Server struct {
	Proto, LocalInet, Port string
	listener               net.Listener
}

func StartServer(params *m.Params) (err error) {
	server = ConfigureServer(params)
	server.listener, err = net.Listen(server.Proto, server.LocalInet+":"+server.Port)

	if err != nil {
		return err
	}
	output := make(chan []byte, params.Queue.MaxQueueN)
	q, err := ConfigureQueue(params)
	if err != nil {
		return err
	}
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn, params *m.Params) {
			buf := make([]byte, params.Queue.MaxMessageL)
			n, err := c.Read(buf)
			if err != nil {
				c.Close()
			}
			output <- buf[:n]
			c.Write(handleMessage(<-output, q))
			c.Close()
		}(conn, params)
	}
	return nil
}

func StopServer() {
	server.listener.Close()
}

func publish(q q.QueueInterface, input []byte) error {
	return q.Push(input)
}

func subscribe(q q.QueueInterface) ([]byte, error) {
	return q.Pop()
}

func handleMessage(message []byte, queue q.QueueInterface) []byte {
	parsed, err := ParseMessage(message)
	if err != nil {
		return []byte("Error parsing the incoming message: " + err.Error())
	}
	switch parsed.Operation {
	case "P":
		decoded, err := base64.StdEncoding.DecodeString(parsed.Payload)
		if err != nil {
			parsed.Error = err
		}
		parsed.Error = publish(queue, decoded)

	case "S":
		resp, err := subscribe(queue)
		parsed.Payload = base64.StdEncoding.EncodeToString(resp)
		parsed.Error = err

	default:
		parsed.Error = errors.New("Error: Operation not implemented")
	}

	if parsed.Error != nil {
		parsed.Confirmed = "N"
	} else {
		parsed.Confirmed = "Y"
		switch queue.(type) {
		case q.Queue:
			syncMessage(message)
		}
	}
	return WriteMessage(parsed)
}

func ConfigureServer(conf *m.Params) *Server {
	var inet, proto, port string
	if conf.Network.Port == "" {
		port = DEFAULT_LISTEN_PORT
	} else {
		port = conf.Network.Port
	}
	if conf.Network.Proto == "" {
		proto = DEFAULT_PROTOCOL
	} else {
		port = conf.Network.Proto
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

func ConfigureQueue(conf *m.Params) (queue q.QueueInterface, err error) {
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
	case m.USE_MEMORY:
		mq := new(q.Queue)
		mq.Init(conf.Queue.MaxQueueC)
		queue = mq
	case m.USE_DATABASE:
		queue = q.DbQueue{}
	case m.USE_FILESYSTEM:
		queue = q.FsQueue{}
	default:
		err = errors.New("Please configure QUEUE_TYPE with 1 (memory), 2 (database) or 3 (filesystem)")
	}
	return queue, err
}
