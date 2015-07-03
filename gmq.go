package main

import (
	"flag"
	c "gmq/communication"
	m "gmq/configuration"
	"log"
	"os"
)

var configfile string
var logger *log.Logger

func main() {

	flag.StringVar(&configfile, "f", "/etc/gmq/gmq.json", "Configuration file")
	flag.Parse()

	config, err := m.ParseConfiguration(configfile)
	if err != nil {
		logger.Fatalf("Could not start server: configuration error\n%T\n%s\n", err, err)
	}
	if logger, err = configureLogger(config); err != nil {
		logger.Printf("Defaulting log to STDOUT because log file is not configured or is unaccessible\n%T\n%s", err, err)
	}
	err = c.StartServer(config)
	if err != nil {
		logger.Printf("Error in GMQ server:\n%T\n%s\n", err, err)
		return
	}
	defer c.StopServer()
}

func configureLogger(p *m.Params) (*log.Logger, error) {
	if p.Log.Path == "" {
		return log.New(os.Stdout, "[GMQ Server] ", log.LstdFlags), nil
	}
	file, err := os.OpenFile(p.Log.Path+p.Log.Name, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return log.New(os.Stdout, "[GMQ Server] ", log.LstdFlags), err
	}
	return log.New(file, "[GMQ Server] ", log.LstdFlags), nil

}
