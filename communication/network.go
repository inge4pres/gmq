package gmqnet

import (
	q "gmq/queue"
	"net"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
	MAX_QUEUES          = 4096
	MAX_MESSAGE_LENGHT  = 40960
)

type Server struct {
	Proto, LocalInet, Port string
}

func (s *Server) StartServer() (chan []byte, error) {
	l, err := net.Listen(s.Proto, s.LocalInet+":"+s.Port)
	if err != nil {
		return nil, err
	}
	defer l.Close()
	output := make(chan []byte, MAX_QUEUES)
	for {
		conn, err := l.Accept()
		if err != nil {
			return nil, err
		}
		go func(c net.Conn) error {
			buf := make([]byte, MAX_MESSAGE_LENGHT)
			if _, err := c.Read(buf); err != nil {
				return err
			}
			output <- buf
			c.Close()
			return nil
		}(conn)
	}
	return output, nil
}

func (s *Server) Publish(q q.QueueInterface, input chan []byte) error {
	return q.Push(<-input)
}

func (s *Server) Subscribe(q q.QueueInterface) ([]byte, error) {
	return q.Pop()
}
