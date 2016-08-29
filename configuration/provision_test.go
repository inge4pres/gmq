package gmqconf

import (
	"fmt"
	"testing"
)

var configfile = "../test/configuration/example_gmq_config.json"

func TestInitServer(t *testing.T) {
	config, err := ParseConfiguration(configfile)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	server, err := InitServer(config)
	netaddr := server.Listener.Addr()
	fmt.Printf("Listner Network: %s", netaddr)
	if netaddr.Network() == "" {
		t.Error("Netowrk Listener not initiated")
		return
	}

}

func TestConfigureQueue(t *testing.T) {
	config, err := ParseConfiguration(configfile)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	_, err = ConfigureQueue(config)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
}
