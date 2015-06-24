package gmq

import (
	"testing"
)

var dsn = "gmq:qm3ssageM3@(sviluppo.mtl.it:3306)/gmq"
var vendor = "mysql"
var dbqueue = &DbQueue{Name: "queue",
	Vendor: vendor,
	Dsn:    dsn,
}

func TestDbQueuePush(t *testing.T) {
	if err := dbqueue.Push(message); err != nil {
		t.Errorf("Error %T %s", err, err)
	}
}

func TestDbQueuePop(t *testing.T) {
	ret, err := dbqueue.Pop()
	if err != nil {
		t.Errorf("Error %T %s", err, err)
	}
	if len(ret) != len(message) {
		t.Errorf("Message pop'd from queue incomplete! \n"+
			"message: %d \n"+
			"returned: %d", len(message), len(ret))
	}
}

func TestDbQueueSequentialPush(t *testing.T) {
	for i := 0; i < 10; i++ {
		if err := dbqueue.Push(message); err != nil {
			t.Errorf("Error %T %s", err, err)
		}
	}
}

func TestDbQueueSequentialPop(t *testing.T) {
	for i := 0; i < 10; i++ {
		ret, err := dbqueue.Pop()
		if err != nil {
			t.Errorf("Error %T %s", err, err)
		}
		if len(ret) != len(message) {
			t.Errorf("Message pop'd from queue incomplete! \n"+
				"message: %d \n"+
				"returned: %d", len(message), len(ret))
		}
	}
}

func TestDbQueueConcurrentPush(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			if err := dbqueue.Push(message); err != nil {
				t.Errorf("Error %T %s", err, err)
			}
		}()
	}
}

func TestDbQueueConcurrentPop(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			ret, err := dbqueue.Pop()
			if err != nil {
				t.Errorf("Error %T %s", err, err)
			}
			if len(ret) != len(message) {
				t.Errorf("Message pop'd from queue incomplete! \n"+
					"message: %d \n"+
					"returned: %d", len(message), len(ret))
			}
		}()
	}
}
