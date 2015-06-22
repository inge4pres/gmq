package gmq

import (
	"testing"
)

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

func TestFsQueueFirstPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
	}
	ret, err := q.Pop()
	// '\n' is a additive byte only used in Push()
	if len(ret) != (len(message) + 1) {
		t.Errorf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
	if err != nil {
		t.Errorf("Error %t %s", err, err)
	}
}

func TestFsQueueSequentialPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
	}
	for i := 0; i < 4; i++ {
		ret, err := q.Pop()
		if len(ret) != (len(message) + 1) {
			t.Errorf("Message pop'd from queue incomplete! \n"+
				"message: %d \n"+
				"returned: %d", len(message), len(ret))
		}
		if err != nil {
			t.Errorf("Error %t %s", err, err)
		}
	}
}

func TestFsQueueConcurrentPop(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_2",
		Path: "../test/fs/",
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
				t.Errorf("Error %t %s", err, err)
			}
		}()
	}
}
