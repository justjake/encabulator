package rekey

import (
	"fmt"
	"github.com/pkg/errors"
	agents "golang.org/x/crypto/ssh/agent"
)

// This data is used to test key signing.
var signedData = []byte("")

// KeyEnsurer is a service object that ensures that the active SSH agent has
// matching keys loaded
type KeyEnsurer struct {
	// If a key matching this predicate is not found, or if signing with it
	// returns an error, we will attempt to re-add the key to the agent using the
	// KeyLoader.
	KeyPredicate func(key *agents.Key) bool
	// A function that attempts to load the key into ssh-agent
	KeyLoader func() error
	agent     agents.Agent
}

// KeyEnsurer creates a new KeyEnsurer.
func New(predicate func(key *agents.Key) bool, loader func() error) *KeyEnsurer {
	return &KeyEnsurer{
		predicate,
		loader,
		nil,
	}
}

// KeyIsLoaded establishes a connection to ssh-agent and queries it for this
// key, returning if the key is present, or if an error occured while querying
// ssh-agent
func (svc *KeyEnsurer) KeyIsLoaded() (bool, error) {
	if (svc.agent) == nil {
		agent, err := ConnectAgent()
		if err != nil {
			return false, errors.Wrap(err, "Connecting to agent")
		}
		svc.agent = agent
	}

	// search keys for
	key, err := FindKey(svc.agent, IsYubikey)
	if err != nil {
		return false, errors.Wrap(err, "Finding key")
	}
	if key == nil {
		return false, nil
	}

	_, err = svc.agent.Sign(key, signedData)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// EnsureLoaded ensures the key is loaded. On success returns nil, otherwise
// returns an error.
func (svc *KeyEnsurer) EnsureLoaded() error {
	loaded, err := svc.KeyIsLoaded()
	if err != nil {
		return errors.Wrap(err, "KeyIsLoaded")
	}
	if loaded {
		return nil
	}

	err = svc.KeyLoader()
	if err != nil {
		return errors.Wrap(err, "KeyLoader")
	}

	loaded, err = svc.KeyIsLoaded()
	if err != nil {
		return errors.Wrap(err, "KeyIsLoaded after loading the key")
	}
	if loaded {
		return nil
	}
	return fmt.Errorf("Key still not loaded after calling loader %#v", svc.KeyLoader)
}

// EnsureRestartingAgent creates a new KeyEnsurer that will restart SSH agent
// before trying to load a key. This works on macOS because ssh-agent runs as a
// user daemon, and SSH_AUTH_SOCK is always updated to the current socket path.
// Your milage may vary on other operating systems.
func EnsureRestartingAgent(svc *KeyEnsurer) *KeyEnsurer {
	copy := &KeyEnsurer{
		svc.KeyPredicate,
		svc.KeyLoader,
		svc.agent,
	}

	loader := func() error {
		if err := KillSSHAgent(); err != nil {
			return errors.Wrap(err, "Killing ssh-agent")
		}
		copy.agent = nil
		return svc.KeyLoader()
	}
	copy.KeyLoader = loader

	return copy
}
