package gmq

import (
	"bytes"
	"io"
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

func getQueueFile(fs *FsQueue) (err error) {
	fs.File, err = os.OpenFile(fs.Path+fs.Name, os.O_RDWR|os.O_CREATE, 0666)
	return
}

func (fs *FsQueue) Push(o []byte) (err error) {
	if err = getQueueFile(fs); err != nil {
		return
	}
	defer fs.File.Close()
	fs.lock.Lock()
	fs.File.Seek(0, os.SEEK_END)
	fs.File.WriteString(string(o) + "\n")
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
	fi, err := fs.File.Stat()
	if err != nil {
		return nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	fs.File.Seek(0, os.SEEK_SET)
	io.Copy(buf, fs.File)
	firstline, err := buf.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil
	}

	fs.File.Seek(0, os.SEEK_SET)
	nw, err := io.Copy(fs.File, buf)
	if err != nil {
		return nil
	}

	fs.File.Truncate(nw)

	ret = []byte(firstline)
	fs.sync()
	fs.lock.Unlock()
	return ret
}

func (fs *FsQueue) sync() {
	fs.File.Sync()
	fs.File.Seek(0, os.SEEK_SET)
}
