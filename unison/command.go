package unison

import (
	"fmt"
	"github.com/justjake/encabulator/command"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

const defaultTerse = false
const macUnisonStateDir = "~/Library/Application Support/Unison"
const linuxUnisonStateDir = "~/.unison"

// macOS doesn't provide stable hostnames ususally,
// so DefaultHostName is used instead.
//
// TODO: learn a host's name once, then write to a file if the host is macOS.
const defaultHostName = "developer-localhost"

var defaultIgnore = []string{}
var defaultIgnoreNot = []string{}

// Unison returns a unison command to use for syncing. An error is returned if
// onlyPaths are not within local.
func Unison(local, remote string, repeat bool, onlyPaths ...string) (*command.T, error) {
	cmd := command.New("unison").
		Flag("root", []string{local, remote}).
		Flag("prefer", local).
		Flag("clientHostName", defaultHostName).
		Flag("ignore", defaultIgnore).
		Flag("ignorenot", defaultIgnoreNot).
		Flag("dumbtty", true).
		Flag("terse", defaultTerse).
		Flag("batch", true)
	if repeat {
		// TOOD: switch to repeat: watch once we're sure we can support it with a
		// go-watchman watcher on all platforms
		cmd.Flag("repeat", "1")
	}
	for i, path := range onlyPaths {
		rel, err := relWithin(local, path)
		if err != nil {
			// TODO: wrap
			return nil, err
		}
		onlyPaths[i] = rel
	}
	cmd.Flag("path", onlyPaths)
	cmd.LongFlagPrefix = &command.Dash
	return cmd, nil
}

// return relative path of target within ancestor base, or return an error if
// target is not in base.
func relWithin(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		// TODO: wrap
		return "", err
	}

	parts := strings.SplitN(rel, "/", 2)
	if len(parts) > 0 {
		if parts[0] == ".." {
			return "", errors.Errorf("Path %q not within root %q", target, base)
		}
	}

	return rel, nil
}

// Remote returns a string specifying a remote host path. Relative paths
// traverse from the users home dir
func Remote(user, host, path string) string {
	return fmt.Sprintf("ssh://%s@%s/%s", user, host, path)
}
