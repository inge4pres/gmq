package gmqnet

import (
	"net"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
	MAX_QUEUES          = 4096
	MAX_MESSSAGE_LENGHT = 10240
)

type Server struct {
	Proto, LocalInet, Port string
}

func (s *Server) Receive() (chan []byte, error) {
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
		go func(c net.Conn) {
			buf := make([]byte, MAX_MESSSAGE_LENGHT)
			if _, err := c.Read(buf); err != nil {
				panic(err)
			}
			output <- buf
			c.Close()
		}(conn)
	}
	return output, nil
}
