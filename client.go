package libgemini

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"
)

type Client struct {
	Config  *tls.Config
	Timeout time.Duration
}

const (
	maxURLSize = 1024
)

// Get will make a request to the url specified by rawURL. To pass in a context
// please use GetWithContext instead.
func (c Client) Get(rawURL string) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.GetWithContext(ctx, rawURL)
}

func (c Client) GetWithContext(ctx context.Context, rawURL string) (Response, error) {
	req, err := NewRequest(rawURL)
	if err != nil {
		return Response{}, err
	}

	return c.DoWithContext(ctx, req)
}

func (c Client) Do(req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.DoWithContext(ctx, req)
}

func (c Client) DoWithContext(ctx context.Context, req Request) (Response, error) {
	c.Config.ServerName = req.u.Hostname()
	d := tls.Dialer{Config: c.Config}

	conn, err := d.DialContext(ctx, "tcp", req.u.Host)
	if err != nil {
		return Response{}, fmt.Errorf("error dialing (%s): %w", req.u.Host, err)
	}
	defer conn.Close()

	fmt.Println("REQ: ", req.u.String())

	if sendErr := SendRequest(conn, req.u.String()); sendErr != nil {
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
