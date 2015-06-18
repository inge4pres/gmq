package gmq

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FsQueue struct {
	lock       sync.RWMutex
	Name, Path string
	File       *os.File
}

type FsPrioQueue struct {
	lock       sync.RWMutex
	Name, Path string
	File       *os.File
}

func createQueueFile(fs *FsQueue) (err error) {
	fs.File, err = os.Create(fs.Path + fs.Name)
	return
}

func getQueueFile(fs *FsQueue) (err error) {
	fs.File, err = os.OpenFile(fs.Path+fs.Name, os.O_RDWR, 660)
	return
}

func appendQueueFile(fs *FsQueue) (err error) {
	fs.File, err = os.OpenFile(fs.Path+fs.Name, os.O_APPEND, 660)
	return
}

func (fs *FsQueue) Push(o []byte) (err error) {
	if err = appendQueueFile(fs); err != nil {
		err = createQueueFile(fs)
	}
	defer fs.File.Close()
	fs.lock.Lock()
	fs.File.Write(o)
	fs.File.WriteString("\n")
	fs.File.Sync()
	fs.lock.Unlock()
	return
}

func (fs *FsQueue) Pop() []byte {
	fs.lock.Lock()
	var ret []byte
	if err := getQueueFile(fs); err != nil {
		return nil
	}
	defer fs.File.Close()
	firstline, err := bufio.NewReader(fs.File).ReadString('\n')
	if err != nil {
		return nil
	}
	ret = []byte(firstline)
	fs.sync()
	fs.lock.Unlock()
	return ret
}

func (fs *FsQueue) sync() {
	var buf []string
	scanner := bufio.NewScanner(fs.File)
	for scanner.Scan() {
		fmt.Println("Reading from file :)")
		line := scanner.Text()
		buf = append(buf, line)
	}
	for s := range buf {
		fs.File.WriteString(buf[s])
	}
	fs.File.Sync()
}
