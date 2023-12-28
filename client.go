package libgemini

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

func NewClient() Client {
	return Client{
		Timeout: defaultTimeout,
		Config: &tls.Config{
			MinVersion:         tls.VersionTLS13,
			InsecureSkipVerify: true,
		},
	}
}

type Client struct {
	Config  *tls.Config
	Timeout time.Duration
}

// Get will call GetWithContext, passing in a context.WithTimeout using c.Timeout.
func (c Client) Get(rawURL string) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.GetWithContext(ctx, rawURL)
}

// GetWithContext will create a Request for the given rawURL and call DoWithContext on it.
func (c Client) GetWithContext(ctx context.Context, rawURL string) (Response, error) {
	req, err := NewRequest(rawURL)
	if err != nil {
		return Response{}, err
	}

	return c.DoWithContext(ctx, req)
}

// Do will call DoWithContext, passing in a context.WithTimeout set to c.Timeout.
// See: DoWithContext for more information.
func (c Client) Do(req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.DoWithContext(ctx, req)
}

// DoWithContext will dial the host, connect to it, finally writing the request on the
// connection.
func (c Client) DoWithContext(ctx context.Context, req Request) (Response, error) {
	c.Config.ServerName = req.u.Hostname()
	d := tls.Dialer{Config: c.Config}

	conn, err := d.DialContext(ctx, "tcp", req.u.Host)
	if err != nil {
		return Response{}, fmt.Errorf("error dialing (%s): %w", req.u.Host, err)
	}
	defer conn.Close()

	if sendErr := req.Write(conn); sendErr != nil {
		return Response{}, fmt.Errorf("error making request: %w", sendErr)
	}

	resp, err := ReadResponse(conn)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
