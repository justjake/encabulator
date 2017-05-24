package commands

import (
	"reflect"
	"testing"
)

func getTestCmd() *Cmd {
	return &Cmd{
		Path: "ruby",
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

func TestCmd_Build(t *testing.T) {
	cmd := getTestCmd()
	out := cmd.Build()

	expectedArgs := []string{"--enable", "foo", "-e", "puts :hello, ARGV", "--bools", "foo", "bar"}

	if out.Path != cmd.Path {
		t.Errorf("Paths not equal: %q != %q", cmd.Path, out.Path)
	}

	if !reflect.DeepEqual(out.Args, expectedArgs) {
		t.Errorf("Args not equal: %q != %q", expectedArgs, out.Args)
	}
}

func TestCmd_String(t *testing.T) {
	getTestCmd().String()
}
