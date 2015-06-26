package gmqnet

import (
	"encoding/json"
)

type Message struct {
	Operation byte   `json:"operation"`
	Queue     string `json:"queue"`
	Payload   []byte `json:"payload"`
	Confirmed byte   `json:"confirmation"`
}

func ParseMessage(in []byte) (*Message, error) {
	m := new(Message)
	return m, json.Unmarshal(in, m)
}

func WriteMessage(m *Message) ([]byte, error) {
	return json.Marshal(m)
}
