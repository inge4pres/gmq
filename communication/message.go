package gmqnet

import (
	"encoding/json"
)

type Message struct {
	Operation byte   `json:"operation"`
	Payload   []byte `json:"payload"`
	Confirmed byte   `json:"confirm"`
}

func ParseMessage(in []byte) (*Message, error) {
	m := new(Message)
	return m, json.Unmarshal(in, m)
}
