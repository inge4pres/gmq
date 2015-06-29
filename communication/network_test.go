package gmqnet

import (
	"testing"
)

func TestStartServer(t *testing.T) {
	s := &Server{Port: DEFAULT_LISTEN_PORT, Proto: "tcp"}
	err := s.StartServer()
	defer s.StopServer()
	if err != nil {
		t.Errorf("Error %T %s", err, err)
	}
}
