package main

import (
	"flag"
	"gmq/communication"
	"gmq/configuration"
	"log"
	"os"
)

var configfile string
var logger *log.Logger

func main() {

	flag.StringVar(&configfile, "f", "/etc/gmq/gmq.json", "Configuration file")
	flag.Parse()

	config, err := gmqconf.ParseConfiguration(configfile)
	if err != nil {
		log.New(os.Stdout, "[GMQ Server]", log.LstdFlags).Fatalf("Could not start server: configuration error\n%T\n%s\n", err, err)
	}
	if logger, err = configureLogger(config); err != nil {
		logger.Printf("Defaulting log to STDOUT because log file is not configured or is unaccessible\n%T\n%s", err, err)
	}

	server, err := gmqconf.InitServer(config)
	if err != nil {
		logger.Fatalf("Could not start TCP server!\n%T %s\n Check configuration json", err)
	}

	if err = gmqnet.HandleConnection(server, config); err != nil {
		logger.Fatalf("Error in GMQ server:\n%T\n%s\n", err, err)
	}
	server.StopServer()
	return
}

func configureLogger(p *gmqconf.Params) (*log.Logger, error) {
	if p.Log.Path == "" {
		return log.New(os.Stdout, "[GMQ Server] ", log.LstdFlags), nil
	}
	file, err := os.OpenFile(p.Log.Path+p.Log.Name, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return log.New(os.Stdout, "[GMQ Server] ", log.LstdFlags), err
	}
	return log.New(file, "[GMQ Server] ", log.LstdFlags), nil

}
