package libgemini

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Timeout time.Duration
	Config  *tls.Config
}

const (
	GeminiPort = 1965
)

func (c Client) Get(requestURL string) (Response, error) {
	u, err := ParseURL(requestURL)
	if err != nil {
		return Response{}, err
	}
	u.Scheme = "gemini"

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	c.Config.ServerName = u.Host
	d := tls.Dialer{Config: c.Config}

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", u.Host, GeminiPort))
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	if err := SendRequest(conn, u.String()); err != nil {
		return Response{}, err
	}

	resp, err := ReadResponse(conn)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

const maxURLSize = 1024

func SendRequest(w io.Writer, u string) error {
	if len(u) > maxURLSize {
		return fmt.Errorf("url length (%d) exceeds max size (%d)", len(u), maxURLSize)
	}

	fmt.Println("URL: ", u)

	n, err := fmt.Fprintf(w, "%s\r\n", u)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if n != len(u)+2 {
		return fmt.Errorf("expected to write %d bytes, wrote %d", len(u)+2, n)
	}

	return nil
}

func ParseURL(requestURL string) (*url.URL, error) {
	u, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "gemini"
	}

	if u.Host == "" {
		s := strings.TrimPrefix(u.String(), u.Scheme+"://")
		indx := strings.Index(s, "/")
		host := s
		if indx > -1 {
			host = s[:indx]
		}
		u.Host = host
	}

	u.Path = strings.TrimPrefix(
		strings.TrimPrefix(requestURL, u.Scheme+"://"),
		u.Host,
	)

	return u, nil
}
