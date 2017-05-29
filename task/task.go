package task

import (
	"bufio"
	"fmt"
	ptylib "github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
)

// Task is a running process inside a PTY. Use Spawn to create a new Task.
// Interact with a task by reading from its Events channel, and writing to its
// Input channel.
type Task struct {
	cmd *exec.Cmd
	pty *os.File
	// Send a byte slice to this channel to write to the process's pty.
	Input chan<- []byte
	// Output will emit Event structs as events (like process output or process
	// termination) occurr.
	Output    <-chan *Event
	splitFunc bufio.SplitFunc
}

func (task *Task) String() string {
	return fmt.Sprintf("%T{%p '%+v'}", task, task, task.cmd.Path)
}

// Kill kills the task. Returns nil if the task is not running.
func (task *Task) Kill() error {
	if task.cmd.Process == nil {
		return nil
	}

	return task.cmd.Process.Kill()
}

// Spawn a Cmd into a PTY. Returns a Task, so that you can communicate with
// the process
func Spawn(cmd *exec.Cmd, splitter bufio.SplitFunc) (*Task, error) {
	pty, err := ptylib.Start(cmd)
	if err != nil {
		// TODO: wrap error.
		return nil, err
	}

	// if we don't make the terminal raw, it will echo all input back to us.
	_, err = terminal.MakeRaw(int(pty.Fd()))
	if err != nil {
		// TODO: wrap error.
		return nil, err
	}

	scanner := bufio.NewScanner(pty)
	scanner.Split(splitter)

	fromProcess := make(chan *Event)
	toProcess := make(chan []byte)

	task := &Task{
		cmd:       cmd,
		pty:       pty,
		splitFunc: splitter,
		Input:     toProcess,
		Output:    fromProcess,
	}

	go emitEvents(task, scanner, fromProcess, toProcess)
	go sendInput(pty, toProcess)

	return task, nil
}

// Respawn spawns a new task with a duplicate of this task's command.
func (task *Task) Respawn() (*Task, error) {
	return Spawn(
		&exec.Cmd{
			Path:        task.cmd.Path,
			Args:        task.cmd.Args,
			Env:         task.cmd.Env,
			Dir:         task.cmd.Dir,
			ExtraFiles:  task.cmd.ExtraFiles,
			SysProcAttr: task.cmd.SysProcAttr,
		},
		task.splitFunc,
	)
}

func emitEvents(task *Task, scanner *bufio.Scanner, out chan<- *Event, closeOnEnd chan []byte) {
	for scanner.Scan() {
		output := &Output{scanner.Text()}
		out <- &Event{task, output}
	}

	if err := scanner.Err(); err != nil {
		// TODO: something more sensible than panic
		// send error on the channel?
		panic(err)
	}

	exit := task.cmd.Wait()
	// TODO: should this be last, and go after the send?
	close(closeOnEnd)

	out <- &Event{task, &Ended{exit}}
	close(out)
}

func sendInput(pty *os.File, in <-chan []byte) {
	for input := range in {
		_, err := pty.Write(input)
		if err != nil {
			// TODO: something more sensible than panic
			// TODO: can we always write? What if the process is dead and the PTY closes??
			panic(err)
		}
	}
}
