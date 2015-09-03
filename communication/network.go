package gmqnet

import (
	"encoding/base64"
	"errors"
	"gmq/configuration"
	"gmq/queue"
	"net"
)

func HandleConnection(server *gmqconf.Server, params *gmqconf.Params) (err error) {

	output := make(chan []byte, params.Queue.MaxQueueN)

	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn, params *gmqconf.Params) {
			buf := make([]byte, params.Queue.MaxMessageL)
			n, err := c.Read(buf)
			if err != nil {
				c.Write([]byte("Error: " + err.Error()))
				c.Close()
			}
			output <- buf[:n]
			c.Write(handleMessage(<-output, gmq.QueueInstance))
			c.Close()
		}(conn, params)
	}
	return nil
}

func handleMessage(message []byte, queues map[string]gmq.QueueInterface) []byte {
	var queue gmq.QueueInterface

	parsed, err := ParseMessage(message)
	if err != nil {
		return []byte("Error parsing the incoming message: " + err.Error())
	}

	if queue, exists := queues[parsed.Queue]; !exists {
		queue.Create(parsed.Queue)
		queues[parsed.Queue] = queue
	}

	switch parsed.Operation {
	case "P":
		decoded, err := base64.StdEncoding.DecodeString(parsed.Payload)
		if err != nil {
			parsed.Error = err
		}
		parsed.Error = publish(parsed.Queue, queue, decoded)

	case "S":
		resp, err := subscribe(parsed.Queue, queue)
		parsed.Payload = base64.StdEncoding.EncodeToString(resp)
		parsed.Error = err

	default:
		parsed.Error = errors.New("Error: Operation not implemented")
	}

	if parsed.Error != nil {
		parsed.Confirmed = "N"
	} else {
		parsed.Confirmed = "Y"
	}
	return WriteMessage(parsed)
}

func publish(qname string, q gmq.QueueInterface, input []byte) error {
	return q.Push(input)

}

func subscribe(qname string, q gmq.QueueInterface) ([]byte, error) {
	return q.Pop()
}
