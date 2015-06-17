package gmq

import (
	"log"
	"os"
	"sync"
)

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "GMQ FS ", log.LstdFlags|log.Ldate|log.Ltime)
}

type FsQueue struct {
	lock       sync.RWMutex
	Name, Path string
	Topic      []string
	File       *os.File
	QObj       map[int][]byte
}

type FsPrioQueue struct {
	Name, Path string
	Topic      []string
	File       *os.File
	QObj       map[int][]byte
}

func (fs *FsQueue) getQueueFile() (f *os.File, err error) {
	//TODO
	return nil, nil
}

func (fs *FsQueue) createQueueFile() (*os.File, error) {
	return os.OpenFile(fs.Path+fs.Name+"_"+string(len(fs.QObj)), os.O_RDWR, 0660)

}

func (fs *FsQueue) Push() (err error) {
	//	fs.File, err = fs.createQueueFile()
	//	fs.File.Write()
	//TODO
	return nil
}

func (fs *FsQueue) Pop() []byte {
	//TODO
	return nil
}

func (fs *FsQueue) sync() {
	//TODO
}
