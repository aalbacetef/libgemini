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

type InMemoryStore struct {
	m  map[string][2]string
	mu sync.Mutex
}

func (store *InMemoryStore) Lookup(address string) (Host, error) {
	data := []byte(address)

	store.mu.Lock()
	v, found := store.m[hashFunc(data)]
	store.mu.Unlock()

	if !found {
		return Host{}, ErrHostNotFound
	}

	return Host{
		Address:     address,
		Fingerprint: v[0],
		Comment:     v[1],
	}, nil
}

func (store *InMemoryStore) Add(h Host) error {
	if _, err := store.Lookup(h.Address); err == nil {
		return ErrHostAlreadyExists
	}

	data := []byte(h.Address)

	store.mu.Lock()
	store.m[hashFunc(data)] = [2]string{h.Fingerprint, h.Comment}
	store.mu.Unlock()

	return nil
}

func (store *InMemoryStore) Delete(address string) error {
	_, err := store.Lookup(address)
	if err != nil {
		return err
	}

	data := []byte(address)

	store.mu.Lock()
	delete(store.m, hashFunc(data))
	store.mu.Unlock()

	return nil
}
