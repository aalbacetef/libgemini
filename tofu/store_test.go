package tofu

import (
	"errors"
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	cases := []testCase{
		{
			"it adds and deletes two hosts",
			[]Host{
				{"domain1.com", "a:b:c", "comment-a"},
				{"domain2.com", "b:c:d", "comment-b"},
			},
		},
	}

	storeGens := []struct {
		fn   func(*testing.T) Store
		name string
	}{
		{
			name: "FileStore",
			fn: func(t *testing.T) Store {
				t.Helper()

				return mkFileStore(t)
			},
		},
		{
			name: "InMemoryStore",
			fn: func(t *testing.T) Store {
				t.Helper()

				return NewInMemoryStore()
			},
		},
	}

	for _, storeGen := range storeGens {
		for _, c := range cases {
			label := storeGen.name + "/" + c.label

			t.Run(label, func(tt *testing.T) {
				store := storeGen.fn(t)
				testStoreDelete(tt, store, c)
				testStoreAddLookup(tt, store, c)
			})
		}
	}
}

type testCase struct {
	label string
	hosts []Host
}

func testStoreAddLookup(t *testing.T, store Store, c testCase) {
	t.Helper()

	for _, h := range c.hosts {
		if err := store.Add(h); err != nil {
			t.Fatalf("error adding to store: %v", err)
		}
	}

	for _, host := range c.hosts {
		found, err := store.Lookup(host.Address)
		if err != nil {
			t.Fatalf("error looking up host: %v", err)
		}

		if found.Address != host.Address {
			t.Fatalf("(address) mismatch: got %s, want %s", found.Address, host.Address)
		}

		if found.Fingerprint != host.Fingerprint {
			t.Fatalf("(fingerprint) mismatch: got %s, want %s", found.Fingerprint, host.Fingerprint)
		}

		if found.Comment != host.Comment {
			t.Fatalf("(comment) mismatch: got %s, want %s", found.Comment, host.Comment)
		}
	}
}

func testStoreDelete(t *testing.T, store Store, c testCase) {
	t.Helper()

	testStoreAddLookup(t, store, c)

	for _, host := range c.hosts {
		if err := store.Delete(host.Address); err != nil {
			t.Fatalf("could not delete host: %v", err)
		}

		_, err := store.Lookup(host.Address)
		if errors.Is(err, ErrHostNotFound) {
			continue
		}

		if err != nil {
			t.Fatalf("error looking up host: %v", err)
		}

		t.Fatalf("matching host found")
	}
}

func mkFileStore(t *testing.T) *FileStore {
	t.Helper()

	file, err := os.CreateTemp("", "tofu-store-test-*")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})

	store, err := NewFileStore(file.Name())
	if err != nil {
		t.Fatalf("could not create store: %v", err)
	}

	return store
}
