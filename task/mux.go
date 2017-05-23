package task

// Mux mixes together one or more event channels. Note that a Mux will never
// close its output channel. Currently 100% lock-free!
//
// TODO: use a WaitGroup to close the mux when all added channels close?
type Mux struct {
	out chan *Event
}

// Add a channel to this mux, piping all events from that channel out.
func (mux *Mux) Add(c chan *Event) <-chan *Event {
	go func() {
		for event := range c {
			mux.out <- event
		}
	}()

	return mux.Out()
}

// Out returns a channel that outputs all the events of channels added with Add.
func (mux *Mux) Out() <-chan *Event {
	return mux.out
}
