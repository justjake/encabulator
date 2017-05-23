package main

import (
	"bufio"
	"fmt"
	"github.com/justjake/encabulator/task"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("tee", "./log.txt")
	tee, err := task.Spawn(cmd, bufio.ScanLines)
	if err != nil {
		panic(err)
	}

	tee.Input <- []byte("hello world\n")

	go func(t *task.Task) {
		time.Sleep(time.Second)
		t.Kill()
	}(tee)

	for event := range tee.Output {
		switch payload := event.Payload.(type) {
		case *task.Output:
			fmt.Println("-> %s", payload.Chunk)
		case *task.Ended:
			fmt.Println("Task ended.")
		}
	}
}
