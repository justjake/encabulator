package commands

import (
	"fmt"
	"os/exec"
)

// LongFlagPrefix is the default prefix appended to multi-character flags in a command.
var LongFlagPrefix = "--"

// ShortFlagPrefix The default prefix applied to single-character flags in a command.
var ShortFlagPrefix = "-"

// Cmd is a factory (with decent erganomics) for creating *exec.Cmd instances.
type Cmd struct {
	Path            string
	Flags           map[string]string
	Args            []string
	LongFlagPrefix  *string
	ShortFlagPrefix *string
}

func (cmd *Cmd) applyDefaults() {
	if cmd.LongFlagPrefix == nil {
		cmd.LongFlagPrefix = &LongFlagPrefix
	}
	if cmd.ShortFlagPrefix == nil {
		cmd.ShortFlagPrefix = &ShortFlagPrefix
	}
}

// Build an *exec.Cmd from this Cmd. This returns a new *exec.Cmd on each call,
// allowing you to use a Cmd as a factory.
func (cmd *Cmd) Build() *exec.Cmd {
	cmd.applyDefaults()

	out := new(exec.Cmd)
	out.Path = cmd.Path
	out.Args = make([]string, 0, len(cmd.Args)+len(cmd.Flags)*2)

	// --foo bar
	// -f bar
	for key, value := range cmd.Flags {
		prefix := cmd.LongFlagPrefix
		if len(key) == 1 {
			prefix = cmd.ShortFlagPrefix
		}
		out.Args = append(out.Args, *prefix+key, value)
	}

	out.Args = append(out.Args, cmd.Args...)

	return out
}

func (cmd *Cmd) String() string {
	return fmt.Sprintf("%T{%s %v %v}",
		cmd, cmd.Path, cmd.Flags, cmd.Args)
}
