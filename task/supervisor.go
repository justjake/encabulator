package task

import (
	"github.com/pkg/errors"
	"time"
)

// Supervisor supervises a task, ressurecting it in a new Task if it ever dies.
// A certain number of failures are allowed per time period
type Supervisor struct {
	// only allow maxFailures in any window period
	maxFailures int
	duration    time.Duration
	window      []time.Time
	//lastFailure *Event
}

// Create a new Supervisor
func MakeSupervisor(maxFailures int, within time.Duration) *Supervisor {
	return &Supervisor{
		maxFailures,
		within,
		make([]time.Time, 0, maxFailures),
		//nil,
	}
}

// Zero resets all counters to zero
func (s *Supervisor) Zero() {
	s.window = make([]time.Time, 0, s.maxFailures)
}

// HandleEvent processes an Event. Returns a Task and an error. The task may be
// an ongoing task, or it could be a new task in the case of a task failure.
func (s *Supervisor) HandleEvent(ev *Event) (*Task, error) {
	switch ended := ev.Payload.(type) {
	default:
		return ev.Task, nil
	case *Ended:
		now := time.Now()

		if len(s.window) < s.maxFailures {
			s.window = append(s.window, now)
			return ev.Task.Respawn()
		}

		// we've had too many errors. Bail!
		if now.Sub(s.window[0]) >= s.duration {
			s.Zero()
			return nil, errors.Errorf("Too many errors: %v within %v. Last: %+v",
				s.maxFailures, s.duration, ended)
		}

		// okay, not quite enough errors. Move everything over
		for i, t := range s.window {
			if i == 0 {
				continue
			}
			// this will get replaced on next iteration, unless this is the last
			// element.
			s.window[i] = now
			s.window[i-1] = t
		}
	}

	return ev.Task.Respawn()
}
