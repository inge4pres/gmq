package gmqnet

import (
	"testing"
)

func TestServerReceive(t *testing.T) {
	s := &Server{Port: DEFAULT_LISTEN_PORT, Proto: "tcp"}
	data, err := s.Receive()
	if err != nil {
		t.Errorf("Error %T %s", err, err)
	}
	close(data)
}
