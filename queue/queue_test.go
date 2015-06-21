package gmq

import "testing"

var (
	message = []byte("TEST Message: you know, for testing..."),
	queue = new(Queue{QName:"testqueue"})


func TestQueueSimplePush(t *testing.T) {
	
}
