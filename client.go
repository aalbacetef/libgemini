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
	Config  *tls.Config
	Timeout time.Duration
}

const (
	GeminiPort = 1965
	maxURLSize = 1024
)

func (c Client) Get(requestURL string) (Response, error) {
	parsedURL, err := ParseURL(requestURL)
	if err != nil {
		return Response{}, err
	}

	parsedURL.Scheme = "gemini"

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	c.Config.ServerName = parsedURL.Host
	d := tls.Dialer{Config: c.Config}

	addr := fmt.Sprintf("%s:%d", parsedURL.Host, GeminiPort)

	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return Response{}, fmt.Errorf("error dialing (%s): %w", addr, err)
	}
	defer conn.Close()

	if sendErr := SendRequest(conn, parsedURL.String()); sendErr != nil {
		return Response{}, fmt.Errorf("error making request: %w", sendErr)
	}

	resp, err := ReadResponse(conn)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func SendRequest(w io.Writer, rawURL string) error {
	if len(rawURL) > maxURLSize {
		return fmt.Errorf("url length (%d) exceeds max size (%d)", len(rawURL), maxURLSize)
	}

	const crlf = "\r\n"
	bytesToWrite := len(rawURL) + len(crlf)

	bytesWritten, err := fmt.Fprintf(w, "%s%s", rawURL, crlf)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if bytesWritten != bytesToWrite {
		return fmt.Errorf("expected to write %d bytes, wrote %d", bytesToWrite, bytesWritten)
	}

	return nil
}

func ParseURL(requestURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL (%s): %w", requestURL, err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "gemini"
	}

	if parsedURL.Host == "" {
		s := strings.TrimPrefix(parsedURL.String(), parsedURL.Scheme+"://")
		indx := strings.Index(s, "/")

		host := s
		if indx > -1 {
			host = s[:indx]
		}

		parsedURL.Host = host
	}

	parsedURL.Path = strings.TrimPrefix(
		strings.TrimPrefix(requestURL, parsedURL.Scheme+"://"),
		parsedURL.Host,
	)

	return parsedURL, nil
}
