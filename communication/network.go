package gmqnet

import (
	"errors"
	q "gmq/queue"
	"net"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
	MAX_QUEUES          = 4096
	MAX_MESSAGE_LENGHT  = 40960
)

var queues q.QueueManager

func init() {
	queues.Obj = make(map[string]*q.Queue, MAX_QUEUES)
}

type Server struct {
	Proto, LocalInet, Port string
	listener               net.Listener
}

func (s *Server) StartServer() (err error) {
	s.listener, err = net.Listen(s.Proto, s.LocalInet+":"+s.Port)
	if err != nil {
		return err
	}
	defer s.listener.Close()
	output := make(chan []byte, MAX_QUEUES)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			buf := make([]byte, MAX_MESSAGE_LENGHT)
			n, err := c.Read(buf)
			if err != nil {
				c.Write([]byte("Error: " + err.Error()))
				c.Close()
			}
			output <- buf[:n]
			c.Write(handleMessage(<-output))
			c.Close()
		}(conn)
	}
	return nil
}

func (s *Server) StopServer() {
	s.listener.Close()
}

func produce(q q.QueueInterface, input []byte) error {
	return q.Push(input)
}

func consume(q q.QueueInterface) ([]byte, error) {
	return q.Pop()
}

func handleMessage(message []byte) []byte {
	parsed, err := ParseMessage(message)
	if err != nil {
		return []byte("Error parsing the message: " + err.Error())
	}
	queue, ok := queues.Obj[parsed.Queue]
	if !ok {
		add := new(q.Queue)
		add.QName = parsed.Queue
		queues.Obj[parsed.Queue] = add
	}
	switch parsed.Operation {
	case "P":
		parsed.Error = produce(queue, parsed.Payload)

	case "S":
		parsed.Payload, parsed.Error = consume(queue)

	default:
		parsed.Error = errors.New("Error: Method not implemented")
	}

	if parsed.Error != nil {
		parsed.Confirmed = "N"
	}
	parsed.Confirmed = "Y"
	return WriteMessage(parsed)
}
