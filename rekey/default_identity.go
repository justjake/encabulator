package rekey

import (
	agents "golang.org/x/crypto/ssh/agent"
	"os"
	"os/exec"
)

// LoadDefaultIdentity runs `ssh-add` to load the user's default identity into
// ssh-agent.
func LoadDefaultIdentity() error {
	cmd := exec.Command("ssh-add")

	// pass along IO so users can type into the loader
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DefaultIdentity returns a KeyEnsurer that ensures that "some" identity is
// loaded. If none is present, the default identity is loaded via `ssh-add`
func DefaultIdentity() *KeyEnsurer {
	return New(
		func(_ *agents.Key) bool { return true },
		LoadDefaultIdentity,
	)
}
