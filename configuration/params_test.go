package gmqconf

import (
	"testing"
)

var config_ok = "../test/configuration/example_gmq_config.json"
var config_err = "../test/configuration/error_gmq_config.json"

func TestParseConfiguration(t *testing.T) {
	_, err := ParseConfiguration(config_ok)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	_, err = ParseConfiguration(config_err)
	if err == nil {
		t.Error("This should have raised an error, JSON malformed on purpose")
		return
	}
}
