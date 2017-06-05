package stager

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	pathlib "path"
)

// URWX is a file mode for an executable, or a directory
const URWX = 0755

// URW is a file mode for a data file
const URW = 0644

// Operation is an operation that can be applied
type Operation interface {
	Apply() error
	String() string
}

// File causes a file to be written
type File struct {
	path string
	from io.Reader
	mode os.FileMode
}

// Apply writes the file to disk
// TODO: check existing file, leave in place if equal
func (f *File) Apply() error {
	if err := ensureParentDir(f.path); err != nil {
		return err
	}

	// using complicated OpenFile instead of Create for that mode support
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, f.mode)
	if err != nil {
		return err
	}
	defer assert(file.Close())

	_, err = io.Copy(file, f.from)
	return err
}

func (f *File) String() string {
	return fmt.Sprintf("File{%s mode %x}", f.path, f.mode)
}

// Link causes a link to be written
type Link struct {
	path string
	to   string
}

// Apply creates a symlink at the given path pointing to the given file
func (l *Link) Apply() error {
	if err := ensureParentDir(l.path); err != nil {
		return err
	}

	// we must remove existing objects to overwrite with a
	// symlink.
	if _, err := os.Lstat(l.path); err == nil {
		err = os.Remove(l.path)
		if err != nil {
			return err
		}
	}

	return os.Symlink(l.to, l.path)
}

func (l *Link) String() string {
	return fmt.Sprintf("Link{%s to %s}", l.path, l.to)
}

// Ops is a group of ops.
type Ops struct {
	Name string
	Ops  []Operation
}

// Apply all the ops
func (ops *Ops) Apply() error {
	return Apply(ops.Ops...)
}

func (ops *Ops) String() string {
	return fmt.Sprintf("Ops{%s}", ops.Name)
}

// Group returns a group of operations.
func Group(name string, ops ...Operation) *Ops {
	res := &Ops{name, []Operation{}}
	for _, op := range ops {
		if op != nil {
			res.Ops = append(res.Ops, op)
		}
	}
	return res
}

// Apply all the operations. Stop at the first error.
func Apply(ops ...Operation) error {
	for _, op := range ops {
		err := op.Apply()
		if err != nil {
			errors.Wrap(err, op.String())
		}
	}

	return nil
}

func ensureParentDir(path string) error {
	return os.MkdirAll(pathlib.Dir(path), 0755)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
