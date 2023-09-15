package libgemini

import (
	"fmt"
	"net/url"
)

const (
	geminiScheme = "gemini"
)

// NewRequest will take a raw url and parse it into a url.URL version.
// It will take care of setting some fields, like the scheme, if they are
// blank.
func NewRequest(u string) (Request, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return Request{}, err
	}

	if uri.Scheme == "" {
		uri.Scheme = geminiScheme
	}

	if uri.Scheme == "http" || uri.Scheme == "https" {
		return Request{}, fmt.Errorf("scheme '%s' not supported", uri.Scheme)
	}

	req := Request{uri}
	if err := req.Valid(); err != nil {
		return req, err
	}

	return req, nil
}

type Request struct {
	u *url.URL
}

func (r Request) String() string {
	return r.u.String()
}

func (r Request) Valid() error {
	return nil
}
