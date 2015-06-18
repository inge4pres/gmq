package gmq

import (
	"testing"
)

var message = []byte("TEST Message: you know, for testing...")

func TestFsQueueFirstPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_1",
		Path: "../test/fs/",
	}
	err := q.Push(message)
	if err != nil {
		t.Errorf("Error: %T %s", err, err)
	}
}

func TestFsQueueSequentialPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
	}
	for i := 0; i < 10; i++ {
		err := q.Push(message)
		if err != nil {
			t.Errorf("Error: %T %s", err, err)
		}
	}
}

//DOESN'T FAIL BUT SIDE-EFFECT IS PRESENT: must check goroutines concurrency pattern
func TestFsQueueConcurrentPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_3",
		Path: "../test/fs/",
	}
	for i := 0; i < 10; i++ {
		go func() {
			err := q.Push(message)
			if err != nil {
				t.Errorf("Error: %T %s", err, err)
			}
		}()
	}
}

func TestFsQueuePop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
	}
	ret := q.Pop()
	// '\n' is a additive byte only used in Push()
	if len(ret) != (len(message) + 1) {
		t.Errorf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
}
