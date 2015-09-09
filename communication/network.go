package gmqnet

import (
	"encoding/base64"
	"errors"
	"gmq/configuration"
	"gmq/queue"
	"net"
)

func HandleConnection(server *gmqconf.Server, params *gmqconf.Params) (err error) {
	//Init singletons
	output := make(chan []byte, params.Queue.MaxQueueN)
	queues := gmq.InitQueueInstance(params.Queue.MaxQueueN)

	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn, params *gmqconf.Params, queues map[string]gmq.QueueInterface) {
			buf := make([]byte, params.Queue.MaxMessageL)
			n, err := c.Read(buf)
			if err != nil {
				c.Write([]byte("Error: " + err.Error()))
				c.Close()
			}
			output <- buf[:n]
			c.Write(handleMessage(params, <-output, queues))
			c.Close()
		}(conn, params, queues)
	}
	return nil
}

func handleMessage(params *gmqconf.Params, message []byte, queues map[string]gmq.QueueInterface) []byte {
	queue, err := gmqconf.ConfigureQueue(params)
	if err != nil {
		return []byte("Error parsing the server coniguration: " + err.Error())
	}
	parsed, err := ParseMessage(message)
	if err != nil {
		return []byte("Error parsing the incoming message: " + err.Error())
	}

	if _, exists := queues[parsed.Queue]; !exists {
		queue.Create(parsed.Queue)
		queues[parsed.Queue] = queue
	} else {
		queue = queues[parsed.Queue]
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
