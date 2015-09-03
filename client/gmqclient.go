package main

import (
	"bytes"
	"flag"
	"gmq/communication"
	"io"
	"log"
	"net"
	"os"
)

var action, qname, payload, protocol, server, logfile, output string
var logger *log.Logger

func main() {

	flag.StringVar(&action, "-a", "", "Action to be performed: Publish / Subscribe")
	flag.StringVar(&qname, "-q", "", "Queue name")
	flag.StringVar(&payload, "-m", "", "Base64 encoded payload to use (only needed with \"-o P\")")
	flag.StringVar(&protocol, "-p", "tcp", "Protocol type: tcp4 or tcp6 (default to systems)")
	flag.StringVar(&logfile, "-l", "", "Log destination: if not set, defaults to STDOUT")
	flag.StringVar(&server, "-s", "localhost:4567", "Server address: (DNS|IP)[:PORT] form ")
	flag.StringVar(&output, "-o", "", "Output file, if empty defaults to STDOUT")
	flag.Parse()

	logger = initLog()

	out := configureOutput()
	logger.Printf("Printing output to %s\n", out)
	resp := callServer(configureCall())
	if _, err := out.Write(resp); err != nil {
		logger.Printf("There was an error printing response from server to %s\n", out)
		logger.Println("Response:\n%s", string(resp))
	}

}

func initLog() *log.Logger {
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return log.New(os.Stdout, "GMQ Client: ", log.LstdFlags)

	}
	return log.New(file, "GMQ Client: ", log.LstdFlags)
}

func configureCall() *gmqnet.Message {
	if action == "" {
		logger.Fatalln("Action cannot be null! Set -a")
	}
	if qname == "" {
		logger.Fatalln("Queue name cannot be null! Set -q")
	}
	return &gmqnet.Message{
		Operation: action,
		Queue:     qname,
		Payload:   payload,
	}
}

func configureOutput() io.Writer {
	if output == "" {
		return os.Stdout
	}
	out, err := os.OpenFile(output, os.O_CREATE, 0660)
	if err != nil {
		logger.Fatalln("Error opening output file!")
	}
	return out
}

func callServer(mex *gmqnet.Message) []byte {
	var resp bytes.Buffer
	conn, err := net.Dial(protocol, server)
	if err != nil {
		logger.Fatalf("Error in call to server:\n%T\n%s\n", err, err)
	}
	if _, err := conn.Write(gmqnet.WriteMessage(mex)); err != nil {
		logger.Fatalf("Error in call to server:\n%T\n%s\n", err, err)
	}
	if _, err := conn.Read(resp.Bytes()); err != nil {
		logger.Fatalf("Error in call to server:\n%T\n%s\n", err, err)
	}
	return resp.Bytes()
}
