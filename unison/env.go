package unison

import (
	"github.com/justjake/encabulator/stager"
	"os/user"
	"path/filepath"
	"runtime"
)

const homedirPath = ".encabulator/sync"

func Manager() *stager.Manager {
	u, _ := user.Current()
	root := filepath.Join(u.HomeDir, homedirPath)
	manifest := OsManifest()
	return stager.NewManager(root, Asset, nil, manifest...)
}

func OsManifest() []string {
	return glibcManifest()
	switch runtime.GOOS {
	default:
		// TOOD: handle other operating systems
		return []string{}
	case "linux":
		switch LibcType() {
		case "musl":
			return []string{
				"linux-musl/unison",
				"linux-musl/unison-fsmonitor",
			}
		default:
			return []string{
				"linux-glibc/unison",
				"linux-glibc/unison-fsmonitor",
			}
		}
	}
}

func glibcManifest() []string {
	return []string{
		"linux-glibc/unison",
		"linux-glibc/unison-fsmonitor",
	}
}
