package tofu

import "sync"

const (
	minMemStoreSize = 20
)

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		m: make(map[string][2]string, minMemStoreSize),
	}
}

// InMemoryStore is a concurrency-safe in-memory store of
// known hosts.
type InMemoryStore struct {
	m  map[string][2]string
	mu sync.Mutex
}

// Add will add a Host to the list of known hosts. If
// the Host has already been added, it will return a
// ErrHostAlreadyExists error instead.
func (store *InMemoryStore) Add(host Host) error {
	if _, err := store.Lookup(host.Address); err == nil {
		return ErrHostAlreadyExists
	}

	data := []byte(host.Address)

	store.mu.Lock()
	store.m[hashFunc(data)] = [2]string{host.Fingerprint, host.Comment}
	store.mu.Unlock()

	return nil
}

// Delete will delete the Host matching address.
// If it has not been added, ErrHostNotFound will be returned.
func (store *InMemoryStore) Delete(address string) error {
	if _, err := store.Lookup(address); err != nil {
		return err
	}

	data := []byte(address)

	store.mu.Lock()
	delete(store.m, hashFunc(data))
	store.mu.Unlock()

	return nil
}

// Lookup will look for the Host matching "address"
// If no Host is found ErrHostNotFound is returned.
func (store *InMemoryStore) Lookup(address string) (Host, error) {
	data := []byte(address)

	store.mu.Lock()
	hostData, found := store.m[hashFunc(data)]
	store.mu.Unlock()

	if !found {
		return Host{}, ErrHostNotFound
	}

	return Host{
		Address:     address,
		Fingerprint: hostData[0],
		Comment:     hostData[1],
	}, nil
}
