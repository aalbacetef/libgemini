package tofu

import (
	"crypto/md5"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
)

type Host struct {
	Address     string
	Fingerprint string
	Comment     string
}

func hashFunc(data []byte) string {
	hash := md5.Sum(data)

	return string(hash[:])
}

// Fingerprint returns the md5 hash of the DER encoded bytes.
func Fingerprint(cert *x509.Certificate) string {
	hash := hashFunc(cert.Raw)
	n := len(hash)
	bdr := &strings.Builder{}

	for _, h := range hash[:n-1] {
		fmt.Fprintf(bdr, "%02X:", h)
	}

	fmt.Fprintf(bdr, "%02X", hash[n-1])

	return bdr.String()
}

var (
	ErrHostAlreadyExists = errors.New("host already exists")
	ErrHostNotFound      = errors.New("host not found")
)

type Store interface {
	// Add will add a host. If it is already known, it is
	// expected implementations will return a ErrHostAlreadyExists.
	Add(h Host) error

	// Delete will delete the host if it found, otherwise
	// it is expected implementations will return a ErrHostNotFound.
	Delete(address string) error

	// Lookup will check if a host is present otherwise it
	// is expected that implementations will return an
	// ErrHostNotFound.
	Lookup(address string) (Host, error)
}

// Verify performs TOFU verification on a Host using the
// provided Store.
// It returns true if the verification passes.
// If a host is not known, it will add it to the known hosts.
func Verify(store Store, host Host) (bool, error) {
	storedHost, err := store.Lookup(host.Address)

	if errors.Is(err, ErrHostNotFound) {
		if addErr := store.Add(host); addErr != nil {
			return false, fmt.Errorf("error verifying, could not add new host: %w", addErr)
		}

		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("error verifying: %w", err)
	}

	return host.Fingerprint == storedHost.Fingerprint, nil
}

// Update will update use the Store's methods to update
// a given host. It deletes the old host, then adds the new
// one.
func Update(store Store, h Host) error {
	if err := store.Delete(h.Address); err != nil {
		return fmt.Errorf("could not delete: %w", err)
	}

	if err := store.Add(h); err != nil {
		return fmt.Errorf("could not add host: %w", err)
	}

	return nil
}
