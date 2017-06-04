/*
Package stager ensures that given a list of dependencies, those dependencies
exist together in a a folder called a "workspace". Each entry in a workspace
links to a file in a "cache", which stores contents addressed by its hash.
*/
package stager

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	pathlib "path"
	"path/filepath"
)

// Source is a provider of assets that may or may not be staged.
type Source func(name string) ([]byte, error)

// T is a Stager instance.
type T struct {
	Cache     string
	Workspace string
	Source    Source
	Fs        afero.Fs
}

// Hash some data in a deterministic way. TODO: use SHA257
func Hash(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// HashFile hashes the file at the given path
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Stage asset named "name" to a path in the workspace named "dest"
func (s *T) Stage(name, dest string) error {
	asset, err := s.Source(name)
	if err != nil {
		// TODO: wrap
		return err
	}
	hash := Hash(asset)
	op1, err := s.cacheUpsert(hash, asset, true)
	if err != nil {
		return err
	}
	op2, err := s.workspaceLink(hash, dest)
	if err != nil {
		// TODO: wrap
		return err
	}
	return Group(fmt.Sprintf("Stage %s %s", name, dest), op1, op2).Apply()
}

func (s *T) cachePath(key string) string {
	middle := len(key) / 2
	prefix := key[:middle]
	suffix := key[middle+1:]
	return pathlib.Join(s.Cache, prefix, suffix)
}

// Ensure that the given data is stored in the cache with the given key
func (s *T) cacheUpsert(hash string, data []byte, careful bool) (Operation, error) {
	path := s.cachePath(hash)

	if FileExists(path) {
		if !careful {
			return nil, nil
		}

		// check that the file at that path has the right hash
		fileHash, err := HashFile(path)
		if err != nil {
			return nil, err
		}
		// done!
		if fileHash == hash {
			return nil, nil
		}
	}

	return &File{path, bytes.NewReader(data), URWX}, nil
}

// Ensure there is a link in the workspace at `dest` pointing to the cached
// file with key `hash`.
func (s *T) workspaceLink(hash, dest string) (Operation, error) {
	path := s.cachePath(hash)
	destpath := pathlib.Join(s.Workspace, dest)
	relpath, err := filepath.Rel(destpath, path)
	if err != nil {
		return nil, err
	}

	loc, err := os.Readlink(dest)
	if err == nil && loc == relpath {
		// link exists and is valid
		return nil, nil
	}

	// ignore other error and try to make things right first
	return &Link{destpath, relpath}, nil
}
