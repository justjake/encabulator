package unison

import (
	"github.com/justjake/encabulator/stager"
	"os/user"
	"path/filepath"
	//"runtime"
)

const homedirPath = ".encabulator/sync"

func Manager() *stager.Manager {
	u, _ := user.Current()
	root := filepath.Join(u.HomeDir, homedirPath)
	manifest := OsManifest()
	return stager.NewManager(root, Asset, nil, manifest...)
}

func OsManifest() []string {
	// TODO: separate manifest for each os
	return []string{
		"linux-glibc/unison",
		"linux-glibc/unison-fsmonitor",
	}
}
