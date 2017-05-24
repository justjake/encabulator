package commands

import (
	"fmt"
	"os/exec"
)

// DashDash --
var DashDash = "--"

// Dash -
var Dash = "-"

// LongFlagPrefix is the default prefix appended to multi-character flags in a command.
var LongFlagPrefix = DashDash

// ShortFlagPrefix The default prefix applied to single-character flags in a command.
var ShortFlagPrefix = Dash

// Cmd is a factory (with decent erganomics) for creating *exec.Cmd structs.
//
type Cmd struct {
	// The prefix of the command.
	//   gitLog := commands.Cmd{First: []string{"git", "log"}}.Flag("short", true)
	First []string
	// Cmd handles several different types when it comes to flag values:
	//   - nil: the flag will not be set.
	//   - string: the flag will be set to the string.
	//   - fmt.Stringer: the flag will be set to `value.String()`.
	//   - []string: the flag will be set once for each string in the slice. For
	//     example, {"add: []string{"etc", "usr", "opt"}} will be built into final
	//     arguments "--add", "etc", "--add", "usr", "--add", "opt"
	//   - bool: if the value is true, the flag is added without an argument.
	//   - any: flag will be added with the value printed using
	//     `fmt.Sprintf("%v", value)`
	Flags map[string]interface{}
	// Args are joined after all the flag options.
	Args []string
	// The prefix attatched to flags if the flags are longer than one character.
	// A nil pointer uses the default of "--".
	LongFlagPrefix *string
	// The prefix attatched to flags if the flags are one character long.
	// A nil pointer uses the default of "-".
	ShortFlagPrefix *string
	// A seperator between the last flag and the first arg.
	// A nil pointer omits the seperator.
	FlagsSeperator *string
}

// Command returns a new *Cmd with the given name.
func Command(path string, more ...string) *Cmd {
	return &Cmd{
		First: append([]string{path}, more...),
	}
}

// Flag sets a flag to a value on this Cmd, and returns the Cmd.
func (cmd *Cmd) Flag(flag string, value interface{}) *Cmd {
	if cmd.Flags == nil {
		cmd.Flags = make(map[string]interface{})
	}

	cmd.Flags[flag] = value
	return cmd
}

// Arg adds all the given strings as positional arguments, and returns the Cmd.
func (cmd *Cmd) Arg(values ...string) *Cmd {
	if cmd.Args == nil {
		cmd.Args = values
	} else {
		cmd.Args = append(cmd.Args, values...)
	}

	return cmd
}

func (cmd *Cmd) applyDefaults() {
	if cmd.LongFlagPrefix == nil {
		cmd.LongFlagPrefix = &LongFlagPrefix
	}
	if cmd.ShortFlagPrefix == nil {
		cmd.ShortFlagPrefix = &ShortFlagPrefix
	}
}

// Slice returns this command's arguments as a slice of strings.
func (cmd *Cmd) Slice() []string {
	cmd.applyDefaults()

	length := len(cmd.First) + len(cmd.Args) + len(cmd.Flags)*2
	out := append(make([]string, 0, length), cmd.First...)

	// --foo bar
	// -f bar
	for key, value := range cmd.Flags {
		prefix := cmd.LongFlagPrefix
		if len(key) == 1 {
			prefix = cmd.ShortFlagPrefix
		}

		flag := *prefix + key

		switch value := value.(type) {
		default:
			out = append(out, fmt.Sprintf("%v", value))
		case nil:
			// pass
		case bool:
			// TODO: handle --foo=false or --no-foo ???
			if value {
				out = append(out, flag)
			}
		case string:
			out = append(out, flag, value)
		case []string:
			for _, v := range value {
				out = append(out, flag, v)
			}
		case fmt.Stringer:
			out = append(out, flag, value.String())
		}
	}

	if cmd.FlagsSeperator != nil {
		out = append(out, *cmd.FlagsSeperator)
	}

	out = append(out, cmd.Args...)
	return out
}

// Build an *exec.Cmd from this Cmd. This returns a new *exec.Cmd on each call,
// allowing you to use a Cmd as a factory.
func (cmd *Cmd) Build() *exec.Cmd {
	argv := cmd.Slice()
	name := argv[0]
	return exec.Command(name, argv[1:]...)
}

// Join prepends this command to another command, mutating and returning the
// second command. This is useful for building multi-level commands.
//
// TODO: copy instead of mutate?
func (cmd *Cmd) Join(inner *Cmd) *Cmd {
	prefix := append(cmd.Slice(), inner.First...)
	inner.First = prefix
	return inner
}

func (cmd *Cmd) String() string {
	return fmt.Sprintf(
		"%T{%q %q %q}",
		cmd, cmd.First, cmd.Flags, cmd.Args,
	)
}
