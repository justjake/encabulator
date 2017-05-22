package logbar

import (
	"io"
	"time"
)

type setLine struct {
	index int
	string
}

/*
Manager manages writing a LogBar to an underlying io.Writer (ususally
os.Stderr) at regular intervals, and also thread-safe updates to the LogBar.

If an error occurs when writing to the underlying writer, Manager panics.
*/
type Manager struct {
	logbar     *LogBar
	writer     io.Writer
	started    bool
	barUpdates chan *setLine
	logUpdates chan []byte
	quit       chan bool
	ticker     *time.Ticker
}

// NewManager returns a new Manager
func NewManager(lb *LogBar, writer io.Writer) *Manager {
	return &Manager{
		lb,
		writer,
		false,
		make(chan *setLine),
		make(chan []byte),
		make(chan bool),
		nil,
	}
}

// SetLine sets a line in the bar
func (m *Manager) SetLine(idx int, line string) {
	m.barUpdates <- &setLine{idx, line}
}

// Write to the log
func (m *Manager) Write(p []byte) (n int, err error) {
	m.logUpdates <- p
	return len(p), nil
}

// Start begins writing the log bar to the writer
func (m *Manager) Start() {
	if m.started {
		return
	}
	m.started = true
	m.ticker = time.NewTicker(time.Second / 60)
	go m.work(m.ticker.C)
}

// Stop stops writing the lob bar to the writer
func (m *Manager) Stop() {
	if !m.started {
		return
	}
	m.quit <- true
	m.ticker.Stop()
	m.ticker = nil
	m.started = false
}

func (m *Manager) work(renderClock <-chan time.Time) {
	m.writeOut()
	var quit bool
	for !quit {
		select {
		case <-renderClock:
			m.writeOut()
		case update := <-m.barUpdates:
			m.logbar.SetLine(update.index, update.string)
		case bytes := <-m.logUpdates:
			m.logbar.Write(bytes)
		case quit = <-m.quit:
			m.writeOut()
			// already set quit
		}
	}
}

func (m *Manager) writeOut() {
	if !m.logbar.damaged {
		return
	}
	if _, err := m.logbar.WriteTo(m.writer); err != nil {
		panic(err)
	}
}
