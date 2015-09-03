package gmqnet

import (
	"gmq/configuration"
	"testing"
	"time"
)

var configfile = "../test/configuration/example_gmq_config.json"

func TestHandleConnection(t *testing.T) {

	config, err := gmqconf.ParseConfiguration(configfile)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}

	server, err := gmqconf.InitServer(config)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}

	go HandleConnection(server, config)

	time.Sleep(time.Second * 60)
	server.StopServer()
}
