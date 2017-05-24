package commands

import (
	"reflect"
	"testing"
)

func getTestCmd() *Cmd {
	return &Cmd{
		Path: "ruby",
		Flags: map[string]string{
			"enable": "foo",
			"e":      "puts :hello, ARGV",
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

	expectedArgs := []string{"--enable", "foo", "-e", "puts :hello, ARGV", "foo", "bar"}

	if out.Path != cmd.Path {
		t.Errorf("Paths not equal: %q != %q", cmd.Path, out.Path)
	}

	if !reflect.DeepEqual(out.Args, expectedArgs) {
		t.Errorf("Args not equal: %q != %q", expectedArgs, out.Args)
	}
}
