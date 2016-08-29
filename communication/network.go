package gmqnet

import (
	"encoding/base64"
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/inge4pres/gmq/configuration"
	"github.com/inge4pres/gmq/queue"
)

// HandleConnection is a function configuring the server to accetp connections
// Params: the pointer to the Server singleton instantiedted by the provisioner, the server configuration parameters
// Returns: any error generated during the connection with clients
//
//Every client connection is handled in a goroutine: the incoming message will be parsed and the operation executed
func HandleConnection(server *gmqconf.Server, params *gmqconf.Params) (err error) {
	//Init singletons
	output := make(chan []byte, params.Queue.MaxQueueN)
	queues := gmq.InitQueueInstance(params.Queue.MaxQueueN)
	errs := make(chan (error), 0)

	conn, err := server.Listener.Accept()
	if err != nil {
		return err
	}
	go func(c net.Conn, params *gmqconf.Params, q map[string]gmq.QueueInterface) {
		buf := make([]byte, params.Queue.MaxMessageL)
		n, err := c.Read(buf)
		if err != nil {
			c.Write([]byte("Error during the connection\n" + err.Error()))
			c.Close()
			errs <- err
		}
		output <- buf[:n]
		c.Write(handleMessage(params, <-output, q))
		c.Close()
	}(conn, params, queues)

	return <-errs
}

func handleMessage(params *gmqconf.Params, message []byte, queues map[string]gmq.QueueInterface) []byte {
	queue, err := gmqconf.ConfigureQueue(params)
	if err != nil {
		mex := &Message{Error: errors.New("Error parsing the server coniguration: " + err.Error()),
			Confirmed: "N",
		}
		return WriteMessage(mex)
	}
	parsed, err := ParseMessage(message)
	if err != nil {
		mex := &Message{Error: errors.New("Error parsing the incoming message: " + err.Error()),
			Confirmed: "N",
		}
		return WriteMessage(mex)
	}

	if !verifyAuth(GenSha512Token(params.Auth.User), GenSha512Token(params.Auth.Password), parsed.Auth.UserTok, parsed.Auth.PwdTok) {
		parsed.Error = errors.New("Authentication failed!")
		parsed.Confirmed = "N"
		return WriteMessage(parsed)
	}

	if _, exists := queues[parsed.Queue]; !exists {
		_, err := queue.Create(parsed.Queue)
		if err != nil {
			parsed.Error = err
		}
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

	case "L":
		parsed.Payload, parsed.Error = list(parsed.Queue, queue)

	case "SYNC":
		//TODO Get the queue in memory/fs and verify that is the same received
		if err := qsync(parsed.Queue, parsed.Payload, queue); err != nil {
			log.Printf("SYNC error: %v\n", err)
		}

	default:
		parsed.Error = errors.New("Error: operation not yet implemented")
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

func list(qname string, q gmq.QueueInterface) (string, error) {
	lenght, err := gmq.GetQueueLength(qname)
	return strconv.Itoa(lenght), err
}

func qsync(qname, qpayloads string, q gmq.QueueInterface) error {
	//TODO
	return nil
}
