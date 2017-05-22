package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/justjake/encabulator/logbar"
	"io"
	"os"
	"time"
	//"path/filepath"
)

func open(path string) (*bufio.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(f), nil
}

func logLines(soFar int, reader io.Reader, lb logbar.Interface) int {
	buffered := bufio.NewReader(reader)
	var line string
	var err error
	for err == nil {
		line, err = buffered.ReadString('\n')
		soFar++
		lb.Write([]byte(line))
		lb.SetLine(1, fmt.Sprintf("Lines Read: %d", soFar))
		time.Sleep(time.Millisecond)
	}
	return soFar
}

func main() {
	flag.Parse()
	paths := flag.Args()

	lb := logbar.NewManager(logbar.New(3), os.Stderr)
	lb.Start()
	lb.SetLine(0, "---------------------------------")
	lb.SetLine(1, "Starting...")
	lb.SetLine(2, paths[0])

	var total = 0
	visitor := func(path string, f os.FileInfo, err error) error {
		lb.SetLine(2, path)
		log, err := open(path)
		if err != nil {
			return err
		}
		total = logLines(total, log, lb)
		return nil
	}
	var err error
	for _, path := range paths {
		visitor(path, nil, nil)
	}
	//err := filepath.Walk(path, visitor)
	lb.Stop()
	fmt.Println()
	fmt.Printf("Encountered %v\n", err)
}
