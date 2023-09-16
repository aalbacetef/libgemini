package libgemini

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	geminiScheme = "gemini"
	geminiPort   = 1965
	schemeDelim  = "://"
)

// NewRequest will take a raw url and parse it into a url.URL version.
// It will take care of setting some fields, like the scheme, if they are
// blank.
func NewRequest(rawURL string) (Request, error) {
	if !strings.Contains(rawURL, schemeDelim) {
		rawURL = geminiScheme + schemeDelim + rawURL
	}

	uri, err := url.Parse(rawURL)
	if err != nil {
		return Request{}, fmt.Errorf("could not parse URL (%s): %w", rawURL, err)
	}

	if uri.Port() == "" && uri.Scheme == geminiScheme {
		uri.Host = fmt.Sprintf("%s:%d", uri.Host, geminiPort)
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
