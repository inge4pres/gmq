package gmqnet

import (
	"io/ioutil"
	"testing"
)

var testfile = "../test/communication/test_message.json"

func TestParseMessage(t *testing.T) {
	read, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	m, err := ParseMessage(read)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
		return
	}
	if m.Payload == "" {
		t.Errorf("Payload is empty!")
		return
	}
}
