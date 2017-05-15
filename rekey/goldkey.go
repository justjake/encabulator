package rekey

import (
	agents "golang.org/x/crypto/ssh/agent"
	"os"
	"regexp"
)

const (
	goldkeyLoc1 = "/usr/local/lib/opensc-pkcs11.so"
	goldkeyLoc2 = "/usr/lib/opensc-pkcs11.so"
)

var (
	goldkeyRegexp = regexp.MustCompile("opensc-pkcs11")
)

// IsGoldkey returns true if the given key is from a Goldkey physical token
func IsGoldkey(key *agents.Key) bool {
	return goldkeyRegexp.MatchString(key.Comment)
}

// LoadGoldkey loads the Golkey opensc-pkcs11 module into ssh-agent
func LoadGoldkey() error {
	if _, err := os.Stat(goldkeyLoc1); err == nil {
		return LoadPKCS11(goldkeyLoc1)
	}

	return LoadPKCS11(goldkeyLoc2)
}

// Yubikey returns a KeyEnsurer that ensures that a Yubikey SSH token is
// availible in ssh-agent.
func Goldkey() *KeyEnsurer {
	return New(IsGoldkey, LoadGoldkey)
}
