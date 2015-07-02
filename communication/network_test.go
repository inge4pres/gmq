package gmqnet

import (
	m "gmq/configuration"
	"testing"
	"time"
)

var configfile = "../test/configuration/example_gmq_config.json"

func TestStartServer(t *testing.T) {
	config, err := m.ParseConfiguration(configfile)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	go func() {
		err = StartServer(config)
		if err != nil {
			t.Errorf("Error %T %s", err, err)
			return
		}
	}()
	time.Sleep(time.Second * 120)
	StopServer()
}
