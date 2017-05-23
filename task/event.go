package task

import (
	"fmt"
)

/*
Event is emitted from a Task. Use a type switch on event.Payload to determine
exactly which event occurred, and inspect its properties.

  func onEvent(event *task.Event) {
    switch payload := event.Payload.(type) {
		case *task.Output:
			fmt.Printf("Output: %s", payload.Chunk)
		case *task.Ended:
			fmt.Printf("Task %v exited with: %v", event.Task, payload.Error)
		}
	}
*/
type Event struct {
	// TODO: should Producer be a Task?
	Task    *Task
	Payload interface{}
}

func (e *Event) String() string {
	return fmt.Sprintf("%T{from %v: %v}", e, e.Task, e.Payload)
}

// Ended is the type of payload indicating the process ended. If the process
// exited 0, Error will be nil. Otherwise it will be an exec.ExitError.
type Ended struct {
	Error error
}

func (p *Ended) String() string {
	return fmt.Sprintf("%T{%v}", p, p.Error)
}

// Output is the type of payload indicating the process output some amount of data.
type Output struct {
	Chunk string
}

func (p *Output) String() string {
	return fmt.Sprintf("%T{%v}", p, p.Chunk)
}
