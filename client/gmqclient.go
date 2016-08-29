package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/inge4pres/gmq/communication"
)

var action, qname, payload, protocol, server, logfile, output, user, password string
var logger *log.Logger

func main() {

	flag.StringVar(&user, "user", "", "Username to authenticate with GMQ")
	flag.StringVar(&password, "password", "", "Password to authenticate to GMQ")
	flag.StringVar(&action, "a", "", "Action to be performed: P = Publish, S = Subscribe, L = List")
	flag.StringVar(&qname, "q", "", "Queue name")
	flag.StringVar(&payload, "m", "", "Base64 encoded payload to use (only needed with \"-a P\")")
	flag.StringVar(&protocol, "p", "tcp", "Protocol type: tcp4 or tcp6 (default to system settings)")
	flag.StringVar(&logfile, "l", "", "Log destination: if not set, defaults to STDOUT")
	flag.StringVar(&server, "s", "localhost:4567", "Server address: (DNS|IP)[:PORT] form ")
	flag.StringVar(&output, "o", "", "Output file, if empty defaults to STDOUT")
	flag.Parse()

	logger = initLog()

	out := configureOutput()
	mex := configureCall()
	resp := callServer(mex)
	if _, err := out.Write(resp); err != nil {
		logger.Printf("There was an error printing response from server to %s\n", out.Name())
		logger.Println("Response:\n%s", string(resp))
	}

}

func initLog() *log.Logger {
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return log.New(os.Stdout, "[GMQ Client] ", log.LstdFlags)

	}
	return log.New(file, "[GMQ Client] ", log.LstdFlags)
}

func configureCall() *gmqnet.Message {
	if user == "" {
		logger.Fatalln("User cannot be null! Set -user")
	}
	if password == "" {
		logger.Fatalln("Password cannot be null! Set -password")
	}
	if action == "" {
		logger.Fatalln("Action cannot be null! Set -a")
	}
	if qname == "" {
		logger.Fatalln("Queue name cannot be null! Set -q")
	}

	return &gmqnet.Message{
		Auth: gmqnet.AuthToken{
			UserTok: gmqnet.GenSha512Token(user),
			PwdTok:  gmqnet.GenSha512Token(password),
		},
		Operation: action,
		Queue:     qname,
		Payload:   payload,
	}
}

func configureOutput() *os.File {
	if output == "" {
		return os.Stdout
	}
	out, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		logger.Fatalln("Error opening output file!")
	}
	return out
}

func callServer(mex *gmqnet.Message) []byte {
	conn, err := net.Dial(protocol, server)
	defer conn.Close()
	if err != nil {
		logger.Fatalf("Error connecting to the server:\n%T\n%s\n", err, err.Error())
	}
	if _, err := conn.Write(gmqnet.WriteMessage(mex)); err != nil {
		logger.Fatalf("Error posting the message to server:\n%T\n%s\n", err, err.Error())
	}
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		logger.Fatalf("Error reading response from server:\n%T\n%s\n", err, err.Error())
	}
	return resp
}
