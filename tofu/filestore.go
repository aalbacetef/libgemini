package tofu

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func NewFileStore(fpath string) (*FileStore, error) {
	if fpath == "" {
		return nil, fmt.Errorf("invalid path provided: '%s'", fpath)
	}

	abspath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: '%s'", fpath)
	}

	store := &FileStore{
		fpath: abspath,
		inmem: NewInMemoryStore(),
	}

	if err := store.readLatest(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return store, store.writeFile()
		}

		return nil, err
	}

	return store, nil
}

// FileStore is implemented as an InMemoryStore that is backed
// by a file. By convention, the public methods can be assumed
// to be concurrency-safe but none of the private methods should
// be assumed to be concurrency-safe as they may not be.
//
// @NOTE: consider if embedding the InMemoryStore and sharing the mutex.
type FileStore struct {
	inmem       *InMemoryStore
	lastUpdated time.Time
	fpath       string
	mu          sync.Mutex
}

func (store *FileStore) getLastUpdated() (time.Time, error) {
	stat, err := os.Stat(store.fpath)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not stat file: %w", err)
	}

	return stat.ModTime(), nil
}

func (store *FileStore) stale() (bool, error) {
	lastModTime, err := store.getLastUpdated()

	return lastModTime.After(store.lastUpdated), err
}

func (store *FileStore) readLatest() error {
	data, err := os.ReadFile(store.fpath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	m, err := parse(string(data))
	if err != nil {
		return fmt.Errorf("could not parse file: %w", err)
	}

	store.inmem.mu.Lock()
	store.inmem.m = m
	store.inmem.mu.Unlock()

	store.lastUpdated = time.Now()

	return nil
}

const minLen = 2

// parse takes in the string contents of the known_hosts
// file and tries to build the known_hosts map.
func parse(s string) (map[string][2]string, error) {
	lines := strings.Split(s, "\n")
	hostMap := make(map[string][2]string, len(lines))

	for k, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		n := len(parts)

		if n < minLen {
			return nil, fmt.Errorf("minimum length %d, got %d", minLen, n)
		}

		hash := parts[0]
		fingerprint := parts[1]

		comment := ""
		if n > minLen {
			comment = parts[2]
		}

		if _, found := hostMap[hash]; found {
			return nil, fmt.Errorf("file is invalid or corrupted, same host found more than once in line: %d", k)
		}

		hostMap[hash] = [2]string{fingerprint, comment}
	}

	return hostMap, nil
}

const fileMode = fs.FileMode(0o600)

// writeFile serializes each Host entry into a
// space-separated list of:
//   - hash(address)
//   - fingerprint
//   - comment
//
// and then saves it to file.
// It assumes that the data has been stored correctly.
func (store *FileStore) writeFile() error {
	store.inmem.mu.Lock()

	k := 0
	n := len(store.inmem.m)
	bdr := &strings.Builder{}

	for hash, val := range store.inmem.m {
		bdr.WriteString(hash + " " + val[0] + " " + val[1])

		if k < n-1 {
			bdr.WriteRune('\n')
		}
	}

	store.inmem.mu.Unlock()

	if err := os.WriteFile(store.fpath, []byte(bdr.String()), fileMode); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	store.lastUpdated = time.Now()

	return nil
}

// Add will add a Host to the list of known hosts. If
// the host has already been added, it will return a
// ErrHostAlreadyExists error instead.
func (store *FileStore) Add(h Host) error {
	if err := store.inmem.Add(h); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	return store.writeFile()
}

// Delete will delete the Host matching address from the
// known hosts.
// If it has not been added, ErrHostNotFound will be returned.
func (store *FileStore) Delete(address string) error {
	if _, err := store.Lookup(address); err != nil {
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if err := store.inmem.Delete(address); err != nil {
		return err
	}

	return store.writeFile()
}

// Lookup will look for the host matching "address"
// in the known_hosts file.
// If no Host is found ErrHostNotFound is returned.
func (store *FileStore) Lookup(address string) (Host, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	stale, err := store.stale()
	if err != nil {
		return Host{}, err
	}

	if stale {
		if err := store.readLatest(); err != nil {
			return Host{}, err
		}
	}

	return store.inmem.Lookup(address)
}
