package logbar

// Interface covers the common ground between a LogBar and a
// Manager.
type Interface interface {
	SetLine(index int, line string)
	Write(p []byte) (n int, err error)
}
