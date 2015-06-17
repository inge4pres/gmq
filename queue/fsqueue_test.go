package gmq

import "testing"

func TestcreateQueueFile(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_1",
		Path: "./test/fs/",
	}
	err := createQueueFile(q)
	if err != nil {
		t.Errorf("Creation of file %s failed", q.Name+q.Path)
	}
}

func TestFsQueueFirstPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_1",
		Path: "./test/fs/",
	}
	message := []byte("TEST Message: you know, for testing...")
	err := q.Push(message)
	if err != nil {
		t.Errorf("Error: %T %s", err, err)
	}
}

func TestFsQueueConcurrentPush(t *testing.T) {
	q := &FsQueue{
		Name: "queue_test_1",
		Path: "./test/fs/",
	}
	message := []byte("TEST Message: you know, for testing...")
	for i := 0; i < 100; i++ {
		go func() {
			err := q.Push(message)
			if err != nil {
				t.Errorf("Error: %T %s", err, err)
			}
		}()
	}
}
