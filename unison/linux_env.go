package unison

import (
	"os/exec"
	"regexp"
)

var isGlibc = regexp.MustCompile(`glibc|GNU libc|EGLIBC`)
var isMusl = regexp.MustCompile(`musl`)

// LddVersion returns a command for running "ldd --version"
func LddVersion() *exec.Cmd {
	return exec.Command("ldd", "--version")
}

// LddType parses some output bytes and returns the libc name
func LddType(output []byte) string {
	if isGlibc.Match(output) {
		return "glibc"
	}
	if isMusl.Match(output) {
		return "musl"
	}
	return ""
}

// LibcType returns the name of Libc the system is currently using, based on
// the output of `ldd --version`. Unknown or error types will return "".
func LibcType() string {
	output, err := LddVersion().Output()
	if err != nil {
		return ""
	}
	return LddType(output)
}
