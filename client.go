package libgemini

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/aalbacetef/libgemini/tofu"
)

func NewClient(userOpts ...OptsFn) (*Client, error) {
	c := &Client{userOpts: userOpts}
	c.refresh()

	return c, nil
}

type Client struct {
	TLSConfig *tls.Config
	userOpts  []OptsFn
	Options
}

func (c *Client) refresh() {
	options := resolveOptions(c.userOpts...)

	c.Options = options
	c.TLSConfig = tlsConfigFromOptions(options)
}

const (
	// NOTE: read section 4.1 of the spec.
	minTLSVersion = tls.VersionTLS12
)

func tlsConfigFromOptions(options Options) *tls.Config {
	store := resolveStore(options.StorePath)

	verifyFn := verifyConn(store)
	if options.Insecure {
		verifyFn = func(tls.ConnectionState) error {
			return nil
		}
	}

	tlsConfig := &tls.Config{
		MinVersion:         minTLSVersion,
		InsecureSkipVerify: true,
		VerifyConnection:   verifyFn,
	}

	return tlsConfig
}

func resolveOptions(userOptions ...OptsFn) Options {
	options := mergeOpts(
		defaultOpts(),
		configOpts(resolveConfigFile()),
		envOpts(),
	)

	for _, fn := range userOptions {
		fn(&options)
	}

	return options
}

type verifyFunc func(tls.ConnectionState) error

func verifyConn(store tofu.Store) verifyFunc {
	return func(state tls.ConnectionState) error {
		peerCerts := state.PeerCertificates
		if len(peerCerts) == 0 {
			return fmt.Errorf("no peer certificates")
		}

		// NOTE: we only care about the leaf.
		leaf := state.PeerCertificates[0]

		host := tofu.Host{
			Address:     state.ServerName,
			Fingerprint: tofu.Fingerprint(leaf),
		}

		valid, err := tofu.Verify(store, host)
		if err != nil {
			return fmt.Errorf("error verifying: %w", err)
		}

		if !valid {
			return fmt.Errorf("invalid certificate")
		}

		return nil
	}
}

// Get will call GetWithContext, passing in a context.WithTimeout using c.Timeout.
func (c *Client) Get(rawURL string) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.GetWithContext(ctx, rawURL)
}

// GetWithContext will create a Request for the given rawURL and call DoWithContext on it.
func (c *Client) GetWithContext(ctx context.Context, rawURL string) (Response, error) {
	req, err := NewRequest(rawURL)
	if err != nil {
		return Response{}, err
	}

	return c.DoWithContext(ctx, req)
}

// Do will call DoWithContext, passing in a context.WithTimeout set to c.Timeout.
// See: DoWithContext for more information.
func (c *Client) Do(req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.DoWithContext(ctx, req)
}

// DoWithContext will dial the host, connect to it, finally writing the request on the
// connection.
func (c *Client) DoWithContext(_ctx context.Context, req Request) (Response, error) {
	ctx, cancel := context.WithCancel(_ctx)
	defer cancel()

	c.refresh()

	traceLogger, err := NewLoggerFromPath(ctx, c.Options.Trace)
	if err != nil {
		return Response{}, err
	}

	traceLogger.Info("Client.DoWithContext", "options", c.Options)

	cfg := c.TLSConfig.Clone()
	cfg.ServerName = req.u.Hostname()

	traceLogger.Info(
		"tls config",
		"ServerName", cfg.ServerName,
		"MinVersion", cfg.MinVersion,
		"bypassing TOFU", c.Options.Insecure,
	)

	headersLogger, err := NewLoggerFromPath(ctx, c.Options.DumpHeaders)
	if err != nil {
		return Response{}, err
	}

	d := tls.Dialer{
		Config: cfg,
	}

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

	headersLogger.Info(
		"Headers",
		"Host", cfg.ServerName,
		"URL", req.String(),
		"Meta", resp.Header.Meta,
		"Status", resp.Header.Status,
	)

	return resp, nil
}
