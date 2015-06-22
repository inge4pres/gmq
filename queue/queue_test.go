package gmq

import "testing"

var message = []byte("TEST Message: you know, for testing...")
var queue = &Queue{QName: "testqueue"}

func TestQueueSimplePush(t *testing.T) {
	err := queue.Push(message)
	if err != nil {
		t.Errorf("Error %T %s", err, err)
	}
	if len(queue.QObj) != 1 {
		t.Error("Push failed!")
	}
}

func TestQueueSimplePop(t *testing.T) {
	ret, err := queue.Pop()
	if err != nil {
		t.Errorf("Error %t %s", err, err)
	}
	if len(ret) != len(message) {
		t.Errorf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
}
