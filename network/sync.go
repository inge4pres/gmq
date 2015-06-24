package gmqsync

import (
	_ "GMQ/queue"
	"net"
)

const (
	DEFAULT_LISTEN_PORT = "4884"
)

type Server struct {
	Addr net.Addr
	Port string
}

func (s *Server) Broadcast(mess []byte) {
}
