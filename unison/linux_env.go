package unison

import (
	"os/exec"
	"regexp"
)

const (
	// Unknown version of LDD not mached below
	Unknown = iota
	// Glibc or compatible (like EGLIBC), like Centos or Ubuntu
	Glibc
	// Musl libc, like Apline linux
	Musl
)

var isGlibc = regexp.MustCompile(`glibc|GNU libc|EGLIBC`)
var isMusl = regexp.MustCompile(`musl`)

// LddVersion returns "ldd --version"
func LddVersion() *exec.Cmd {
	return exec.Command("ldd", "--version")
}

// Parse the Libc provider given LDD output
func LddType(output string) int {
	if isGlibc.MatchString(output) {
		return Glibc
	}
	if isMusl.MatchString(output) {
		return Musl
	}
	return Unknown
}
