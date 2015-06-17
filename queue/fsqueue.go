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
	fs.File, err = os.OpenFile(fs.File.Name(), os.O_RDONLY, 660)
	return
}

func appendQueueFile(fs *FsQueue) (err error) {
	fs.File, err = os.OpenFile(fs.Path+fs.Name, os.O_APPEND, 660)
	return
}

func (fs *FsQueue) Push(o []byte) (err error) {
	if err = appendQueueFile(fs); err != nil {
		fmt.Println("Creating an empty file")
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
		panic(err)
	}
	defer fs.File.Close()
	firstline, err := bufio.NewReader(fs.File).ReadString('\n')
	if err != nil {
		panic(err)
	}
	ret = []byte(firstline)
	fs.sync()
	fs.lock.Unlock()
	return ret
}

func (fs *FsQueue) sync() {
	var lines []string
	scanner := bufio.NewScanner(fs.File)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	createQueueFile(fs)
	for n := 1; n < len(lines); n++ {
		fs.File.WriteString(lines[n])
		fs.File.WriteString("\n")
	}
	fs.File.Sync()
}
