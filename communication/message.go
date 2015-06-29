package gmqnet

import (
	"encoding/json"
)

type Message struct {
	Operation string `json:"operation"`
	Queue     string `json:"queue"`
	Payload   []byte `json:"payload"`
	Confirmed string `json:"confirmation"`
	Error     error  `json:"error"`
}

func ParseMessage(in []byte) (*Message, error) {
	m := new(Message)
	return m, json.Unmarshal(in, m)
}

func WriteMessage(m *Message) []byte {
	resp, err := json.Marshal(m)
	if err != nil {
		resp = []byte("ERROR")
	}
	return resp
}
