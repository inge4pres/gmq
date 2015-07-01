package gmqnet

import (
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	s := &Server{Port: DEFAULT_LISTEN_PORT, Proto: "tcp"}
	go func(s *Server) {
		err := s.StartServer()
		if err != nil {
			t.Errorf("Error %T %s", err, err)
		}
	}(s)
	time.Sleep(time.Second * 120)
	s.StopServer()
}
