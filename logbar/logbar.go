/*
Package logbar implements a terminal UI component called a LogBar.

The top part, the "log", is a feed of newline-separated log messages.

The bottom part, the "bar", is a persistent status display that can be filled
with any sorts of characters.
*/
package logbar

import (
	"bytes"
	"github.com/justjake/encabulator/cursor"
	"io"
)

// LogBar is a terminal component. Set its persistent lines with SetLine(),
// write to it with fmt.Fprint, fmt.Fprintf, and fmt.Fprintln.
type LogBar struct {
	// text
	bar []string
	// buffers of new log lines that need to be written to the screen
	logLines *bytes.Buffer
	// true when an update needs to be written
	damaged bool
	// when we write out the LogBar, we'll record the height here to know how
	// many rows up to go
	lastHeight int
}

// New returns a new LogBar
func New(height int) *LogBar {
	return &LogBar{
		bar:        make([]string, height),
		logLines:   new(bytes.Buffer),
		damaged:    false,
		lastHeight: height,
	}
}

// BarHeight returns the height of the bar, in lines
func (lb *LogBar) BarHeight() int {
	return len(lb.bar)
}

// SetLine updates given line in the Bar.
func (lb *LogBar) SetLine(row int, line string) {
	lb.damaged = true

	// base case: ez just set the line
	if row < len(lb.bar) {
		lb.bar[row] = line
		return
	}

	// row is beyoned end of the bar: extend the bar
	extra := (row + 1) - len(lb.bar)
	newLines := make([]string, extra)
	newLines[extra-1] = line
	lb.bar = append(lb.bar, newLines...)
}

// Write to the LogBar's log
func (lb *LogBar) Write(p []byte) (n int, err error) {
	lb.damaged = true
	return lb.logLines.Write(p)
}

// Render a Buffer that can be written to an io.Writer to render this component
// to a terminal.
func (lb *LogBar) Render() *bytes.Buffer {
	var b bytes.Buffer

	// go to the first line of the bar, and clear it, so that we can place our
	// none-or-more log lines, and then redraw our bar.
	b.WriteString(cursor.MoveToHorizontal(0))
	if ups := lb.lastHeight - 1; ups > 0 {
		b.WriteString(cursor.MoveUp(ups))
	}
	b.WriteString(cursor.ClearEntireLine())

	// emit the log lines
	if _, err := lb.logLines.WriteTo(&b); err != nil {
		panic(err)
	}

	// write each line, save the last
	last := len(lb.bar) - 1
	for i, line := range lb.bar {
		b.WriteString(line)
		b.WriteString(cursor.ClearLineRight())
		if i != last {
			b.WriteRune('\n')
		}
	}

	return &b
}

// Clean marks this LogBar as clean, and remembers its height. Call this once
// you write the LogBar out to a writer
func (lb *LogBar) Clean() {
	lb.logLines.Reset()
	lb.lastHeight = len(lb.bar)
	lb.damaged = false
}

// WriteTo writes this to a writer, and cleans it
func (lb *LogBar) WriteTo(wr io.Writer) (int64, error) {
	num, err := lb.Render().WriteTo(wr)
	if err != nil {
		return num, err
	}
	lb.Clean()
	return num, nil
}
