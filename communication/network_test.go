package gmqnet

import (
	"testing"
	"time"

	"github.com/inge4pres/gmq/configuration"
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

	time.Sleep(time.Second * 10)
	server.StopServer()
}
