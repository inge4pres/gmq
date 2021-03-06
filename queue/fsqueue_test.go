package gmq

import (
	"sync"
	"testing"
)

func TestFsQueueFirstPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_1",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	err := q.Push(message)
	if err != nil {
		t.Errorf("Error: %T %s", err, err.Error())
	}
}

func TestFsQueueSequentialPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	for i := 0; i < 10; i++ {
		err := q.Push(message)
		if err != nil {
			t.Errorf("Error: %T %s", err, err.Error())
		}
	}
}

//DOESN'T FAIL BUT SIDE-EFFECT IS PRESENT: must check goroutines concurrency pattern
func TestFsQueueConcurrentPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_3",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	for i := 0; i < 10; i++ {
		go func() {
			err := q.Push(message)
			if err != nil {
				t.Errorf("Error: %T %s", err, err.Error())
			}
		}()
	}
}

func TestFsQueueFirstPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	ret, err := q.Pop()
	// '\n' is an additive byte only used in Push()
	if len(ret) != (len(message) + 1) {
		t.Errorf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
	if err != nil {
		t.Errorf("Error %T %s", err, err.Error())
	}
}

func TestFsQueueSequentialPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	for i := 0; i < 4; i++ {
		ret, err := q.Pop()
		if len(ret) != (len(message) + 1) {
			t.Errorf("Message pop'd from queue incomplete! \n"+
				"message: %d \n"+
				"returned: %d", len(message), len(ret))
		}
		if err != nil {
			t.Errorf("Error %T %s", err, err)
		}
	}
}

func TestFsQueueConcurrentPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
		lock: &sync.RWMutex{},
	}
	for i := 0; i < 5; i++ {
		go func() {
			ret, err := q.Pop()
			if len(ret) != (len(message) + 1) {
				t.Errorf("Message pop'd from queue incomplete! \n"+
					"message: %d \n"+
					"returned: %d", len(message), len(ret))
			}
			if err != nil {
				t.Errorf("Error %T %s", err, err)
			}
		}()
	}
}
