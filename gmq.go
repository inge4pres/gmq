package main

import (
	_ "errors"
	"flag"
	"fmt"
	q "gmq/queue"
	"os"
)

var configfile string
var message = []byte("TEST Message: you know, for testing...")

func main() {

	flag.StringVar(&configfile, "f", "gmq.conf", "Location of configuation file")
	flag.Parse()

	//	if err := parseConfigFile(configfile); err != nil {
	//		panic(errors.New("Unable to parse configuration file %s"))
	//	}
	qu := &q.FsQueue{
		Name: "queue_test_2",
		Path: "./test/fs/",
	}
	for i := 0; i < 10; i++ {
		err := qu.Push(message)
		if err != nil {
			fmt.Printf("Error: %T %s", err, err)
		}
	}

	//SECOND
	ret := qu.Pop()
	// '\n' is a additive byte only used in Push()
	if len(ret) != (len(message) + 1) {
		fmt.Printf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
}

func parseConfigFile(f string) error {
	//TODO
	_, err := os.Open(f)
	return err
}
