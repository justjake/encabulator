package main

import (
	"fmt"
	"github.com/justjake/unison-wrapper/rekey"
	"os"
)

// TODO: add help
// TODO: add support for goldkey
func main() {
	yk := rekey.EnsureRestartingAgent(rekey.Yubikey())
	err := yk.EnsureLoaded()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
