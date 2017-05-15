// Package rekey contains KeyEnsurer, a type that ensures that certain SSH keys
// are availible in your ssh-agent.
//
// Specific support is provided for loading Goldkey and Yubikey token keys on
// macOS.
package rekey

import (
	"fmt"
	"github.com/pkg/errors"
	agents "golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"os/exec"
)

const (
	// SSHAuthSock contains the name of the environment variable used to locate
	// the unix socket to connect to an SSH agent.
	SSHAuthSock = "SSH_AUTH_SOCK"
	// AgentLifetime contains a string representing the duration that keys added
	// to the agent should be available for. This is used in the PKCS11KeyLoader.
	AgentLifetime = "14400"
)

// ConnectAgent returns a new ssh/agent connected to the system default SSH
// agent over SSH_AUTH_SOCK, or an error if a connection cannot be established.
func ConnectAgent() (agents.Agent, error) {
	sockPath, sockSet := os.LookupEnv(SSHAuthSock)
	if !sockSet {
		return nil, fmt.Errorf("Can't connect to SSH Agent because %s is unset", SSHAuthSock)
	}
	sock, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, err
	}
	return agents.NewClient(sock), nil
}

// LoadPKCS11 loads the given PKCS11 path into ssh-agent using `ssh-add`.
func LoadPKCS11(path string) error {
	cmd := exec.Command(
		"ssh-add",
		"-s", path,
		"-t", AgentLifetime,
	)

	// pass along IO so users can type into the loader
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// FindKey searches ssh-agent for the first key that matches the given predicate.
func FindKey(agent agents.Agent, predicate func(*agents.Key) bool) (*agents.Key, error) {
	keys, err := agent.List()
	if err != nil {
		return nil, errors.Wrap(err, "Could not list ssh-agent's keys")
	}

	for _, k := range keys {
		if predicate(k) {
			return k, nil
		}
	}

	return nil, nil
}

func KillSSHAgent() error {
	return exec.Command("killall", "ssh-agent").Run()
}
