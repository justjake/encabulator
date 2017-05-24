package commands

import (
	"fmt"
	"github.com/justjake/encabulator/assert"
	"testing"
)

func getTestCmd() *Cmd {
	return &Cmd{
		First: []string{"ruby"},
		Flags: map[string]interface{}{
			"enable": "foo",
			"e":      "puts :hello, ARGV",
			"bools":  true,
		},
		Args: []string{
			"foo",
			"bar",
		},
	}
}

func TestCommand(t *testing.T) {
	assert.Equal(t, *Command("foo"), Cmd{First: []string{"foo"}})
	assert.Equal(t, *Command("foo", "bar"), Cmd{First: []string{"foo", "bar"}})
}

func TestCmd_Arg(t *testing.T) {
	assert.Equal(t,
		*Command("foo").Arg("bar"),
		Cmd{First: []string{"foo"}, Args: []string{"bar"}},
	)
	assert.Equal(t,
		*Command("foo").Arg("bar", "baz"),
		Cmd{First: []string{"foo"}, Args: []string{"bar", "baz"}},
	)
}

func TestCmd_Flag(t *testing.T) {
	assert.Equal(t,
		*Command("foo").Flag("bar", true),
		Cmd{First: []string{"foo"}, Flags: map[string]interface{}{"bar": true}},
	)
}

func TestCmd_Slice(t *testing.T) {
	// test regular flags and args
	assert.Equal(t,
		getTestCmd().Slice(),
		[]string{"ruby", "--enable", "foo", "-e", "puts :hello, ARGV", "--bools", "foo", "bar"},
	)

	// test array flags
	assert.Equal(t,
		Command("foo").Flag("arg", []string{"one", "two"}).Slice(),
		[]string{"foo", "--arg", "one", "--arg", "two"},
	)

	// test nil array flag
	assert.Equal(t,
		Command("foo").Flag("destroy-universe", nil).Slice(),
		[]string{"foo"},
	)

	// test flag seperator
	c2 := Command("foo").Flag("what", true).Arg("one", "two")
	c2.FlagsSeperator = &DashDash
	assert.Equal(t,
		c2.Slice(),
		[]string{"foo", "--what", "--", "one", "two"},
	)
}

func TestCmd_Join(t *testing.T) {
	joined := (&Cmd{First: []string{"foo"}, Flags: map[string]interface{}{"bar": true}}).Join(Command("doggo"))
	assert.Equal(t, *joined, Cmd{First: []string{"foo", "--bar", "doggo"}})
}

func TestCmd_Build(t *testing.T) {
	cmd := getTestCmd()
	out := cmd.Build()

	expectedArgs := []string{"ruby", "--enable", "foo", "-e", "puts :hello, ARGV", "--bools", "foo", "bar"}

	assert.Equal(t, out.Args, expectedArgs)
}

func TestCmd_String(t *testing.T) {
	getTestCmd().String()
}

// Demonstrates joining an outer command ("git") into an inner command ("log")
// to produce a full multi-leveled command. This techinique is also useful for
// building SSH commands.
func ExampleCmd_Join() {
	cmd := Command("git").Flag("C", "~/src/alt").Join(Command("log")).Flag("short", true).Arg("./foo/bar")

	fmt.Printf("%#v\n", cmd.Slice())
	// Output: []string{"git", "-C", "~/src/alt", "log", "--short", "./foo/bar"}
}
