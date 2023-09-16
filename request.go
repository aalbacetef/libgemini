package libgemini

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

const (
	geminiScheme = "gemini"
	geminiPort   = 1965
	schemeDelim  = "://"
)

// NewRequest will take a raw url and parse it into a url.URL version.
// It will take care of setting some fields, like the scheme and port, if they are
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

// Request is a simple struct which wraps around a url.URL, providing a few methods around it.
type Request struct {
	u *url.URL
}

func (r Request) String() string {
	return r.u.String()
}

const maxRequestSize = 1024

func (r Request) Valid() error {
	n := len(r.u.String())
	if n > maxRequestSize {
		return fmt.Errorf("max request size of %d bytes exceeded, have %d bytes", maxRequestSize, n)
	}

	if r.u.Hostname() == "" {
		return fmt.Errorf("no hostname")
	}

	return nil
}

const (
	CRLF = "\r\n"
)

// Write will take any io.Writer and write the request on it
// in the format of: <URL>\r\n
// returning errors if the bytes written did not match the expected amount.
func (r Request) Write(w io.Writer) error {
	bytesToWrite := []byte(r.String() + CRLF)
	n := len(bytesToWrite)

	wrote, err := w.Write(bytesToWrite)
	if err != nil {
		return fmt.Errorf("error writing request: %w", err)
	}

	if wrote != n {
		return fmt.Errorf("expected to write %d bytes, wrote %d bytes instead", n, wrote)
	}

	return nil
}
