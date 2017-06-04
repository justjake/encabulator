package stager

import (
	"crypto/sha256"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"path/filepath"
	"sort"
)

// Manager manages staging a manifest of dependencies using a stager. Manager's
// main job is to infer the name of the workspace directory based on the SHAs
// of the dependencies in the manifest.
type Manager struct {
	path     string
	source   Source
	fs       afero.Fs
	manifest []string
	stager   *T
}

func (m *Manager) sumManifest() (string, error) {
	hash := sha256.New()
	for _, name := range m.manifest {
		asset, err := m.source(name)
		if err != nil {
			return "", err
		}
		hash.Write(asset)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (m *Manager) newStager() (*T, error) {
	sum, err := m.sumManifest()
	if err != nil {
		return nil, err
	}

	return &T{
		Workspace: filepath.Join(m.path, "workspace", sum),
		Cache:     filepath.Join(m.path, "cache"),
		Source:    m.source,
		Fs:        m.fs,
	}, nil
}

// NewManager returns a new manager. Creating a manager hashes all the assets
// in the manifest, so creating a manager is a) expensive, and b) may return an
// error.
func NewManager(path string, source Source, fs afero.Fs, manifest ...string) *Manager {
	sort.Strings(manifest)
	return &Manager{path, source, fs, manifest, nil}
}

// Stage a file into the managed workspace
func (m *Manager) Stage(name, dest string) error {
	if m.stager == nil {
		stager, err := m.newStager()
		if err != nil {
			return err
		}
		m.stager = stager
	}

	i := sort.SearchStrings(m.manifest, name)
	if i < len(m.manifest) && m.manifest[i] == name {
		return m.stager.Stage(name, dest)
	}

	return errors.Errorf("Requested asset %q not in manifest", name)
}

// Workspace returns the path to the managed workspace
func (m *Manager) Workspace() string {
	return m.stager.Workspace
}
