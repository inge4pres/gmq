package gmq

import (
	"bufio"
	"bytes"
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
	fs.File, err = os.Open(fs.Path + fs.Name)
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
	//	var lines []string
	//	scanner := bufio.NewReader(fs.File)
	//	for line, err := scanner.ReadString('\n') {
	//		lines = append(lines, line)
	//		fmt.Printf("Reading %d lines from buffer", len(lines))
	//	}
	//	os.Remove(fs.File.Name())
	//	fs.File = nil
	//	if err := createQueueFile(fs); err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("Created new FILE and now starting write on it")
	//	for n := 1; n < len(lines); n++ {
	//		fs.File.WriteString(lines[n])
	//		fmt.Printf("Line %d written on FILE %s", n, fs.File.Name())
	//		fs.File.WriteString("\n")
	//	}
	//	fmt.Println("Closing FILE")
	var current, final bytes.Buffer
	read, _ := fs.File.Read(current.Bytes())
	fmt.Printf("Read %d bytes", read)
	off := bytes.IndexByte(current.Bytes(), '\n') + 1
	n, _ := fs.File.ReadAt(final.Bytes(), int64(off))
	fmt.Printf("Read %d bytes after offset", n)
	fs.File.Write(final.Bytes())
	fmt.Printf("BUFFER LENGHTS: current %d  ,  final %d", len(current.Bytes()), len(final.Bytes()))
	fs.File.Sync()
}
