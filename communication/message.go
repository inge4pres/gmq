package gmqnet

import (
	"encoding/json"
)

/*
Message:
A JSON file with attributes

Operation: "P" (publish, post a message) or "S" (subscribe, retrieve a message)
Queue: the queue name to be used
Payload: in case Operation is "P", the base64 encoded message to be stored in the queue
Confirmed: contains the result of submitted operation, can be "Y" or "N"
Error: strin representation of the error, if none empty
*/
type Message struct {
	Operation string `json:"operation"`
	Queue     string `json:"queue"`
	Payload   string `json:"payload"`
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
		resp = []byte("ERROR: writing the JSON message failed\n" + err.Error())
	}
	return resp
}
