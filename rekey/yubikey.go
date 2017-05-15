package rekey

import (
	agents "golang.org/x/crypto/ssh/agent"
	"regexp"
)

var (
	ykRegexp = regexp.MustCompile("libykcs11|yubico")
)

// IsYubikey returns true if the given key is from a Yubikey physical token
func IsYubikey(key *agents.Key) bool {
	return ykRegexp.MatchString(key.Comment)
}

// LoadYubikey loads the Yubikey PKCS11 lib
func LoadYubikey() error {
	return LoadPKCS11("/usr/local/lib/libykcs11.dylib")
}

// Yubikey returns a KeyEnsurer that ensures that a Yubikey SSH token is
// availible in ssh-agent.
func Yubikey() *KeyEnsurer {
	return New(IsYubikey, LoadYubikey)
}
