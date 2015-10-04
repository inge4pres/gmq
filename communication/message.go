package gmqnet

import (
	"encoding/json"
)

/*
Message is the type describing the information exchange between gmqserver and gmqclient

Operation: "P" (publish, post a message) or "S" (subscribe, retrieve a message)
Queue: the queue name to be used
Payload: in case Operation is "P", the base64 encoded message to be stored in the queue
Confirmed: contains the result of submitted operation, can be "Y" or "N"
Error: strin representation of the error, if none empty
*/
type Message struct {
	Operation string    `json:"operation"`
	Queue     string    `json:"queue"`
	Payload   string    `json:"payload"`
	Confirmed string    `json:"confirmation"`
	Error     error     `json:"error"`
	Auth      AuthToken `json:"auth"`
}

//Authentication hashes to acces the server
type AuthToken struct {
	UserTok string `json:"user"`
	PwdTok  string `json:"password"`
}

// Parses a JSON byte array trying to convert it to a Message
// Params: the input byte array
// Returns: the pointer to the Mesage type, an error if parsing fails
func ParseMessage(in []byte) (*Message, error) {
	m := new(Message)
	return m, json.Unmarshal(in, m)
}

// Write a Message type into a byte array
// Params: the pointer to the Message type
// Returns: the JSON byte array
func WriteMessage(m *Message) []byte {
	resp, err := json.Marshal(m)
	if err != nil {
		resp = []byte("ERROR: writing the JSON message failed\n" + err.Error())
	}
	return resp
}
