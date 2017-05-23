/*
Package task provides building blocks for spawning processes and communicating
with them via channels.

Here's an example of communicating with a simple process and reading its output.

	package main

	import (
		// needed for ScanLines split function
		"bufio"
		"fmt"
		"github.com/justjake/encabulator/task"
		"os/exec"
		// needed to for sleep. Normally not required for working with tasks.
		"time"
	)

	func main() {
		// `tee` reads from STDIN, writes each line to the given file,
	  // then echos to STDOUT.
		cmd := exec.Command("tee", "./log.txt")

		// Spawn `tee` as a task, so we can communicate with it.
		tee, err := task.Spawn(cmd, bufio.ScanLines)
		if err != nil {
			panic(err)
		}

		// write "hello world\n" to tee's STDIN.
		tee.Input <- []byte("hello world\n")

		// start a goroutine to kill tee after a second. This is needed as tee will
		// run forever.
		go func(t *task.Task) {
			time.Sleep(time.Second)
			t.Kill()
		}(tee)

		// Loop over tee's output channel, handling each event.
		// This channel will close once the task's command exits.
		for event := range tee.Output {
			// use a type switch to handle the different sorts of events.
			switch payload := event.Payload.(type) {
			case *task.Output:
				fmt.Println("-> %s", payload.Chunk)
			case *task.Ended:
				fmt.Println("Task ended.")
			}
		}
	}
*/
package task
