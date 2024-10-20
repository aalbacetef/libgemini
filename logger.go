package libgemini

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type FileHandler struct {
	wc io.WriteCloser
	slog.Handler
}

func NewFileHandler(ctx context.Context, wc io.WriteCloser) FileHandler {
	go func() {
		<-ctx.Done()
		wc.Close()
	}()

	return FileHandler{
		wc:      wc,
		Handler: slog.NewJSONHandler(wc, nil),
	}
}

type NoopHandler struct{}

func (NoopHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (NoopHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (NoopHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return NoopHandler{}
}

func (NoopHandler) WithGroup(string) slog.Handler {
	return NoopHandler{}
}

type NoopCloser struct {
	io.Writer
}

func (NoopCloser) Close() error {
	return nil
}

const (
	fileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	filePerms = fs.FileMode(0o644)
	StdOutKey = ":stdout:"
	StdErrKey = ":stderr:"
)

// NewLoggerFromPath will take in a path and if it's empty it will
// return a logger using a NoopHandler, otherwise it will try to
// initialize a FileHandler based on the filepath. It is
// expected to be used alongside a context.
// Use :stderr: or :stdout: to use StdErr or StdOut as outputs.
func NewLoggerFromPath(ctx context.Context, fpath string) (*slog.Logger, error) {
	if fpath == "" {
		return slog.New(NoopHandler{}), nil
	}

	if fpath == StdOutKey {
		return slog.New(NewFileHandler(ctx, NoopCloser{os.Stdout})), nil
	}

	if fpath == StdErrKey {
		return slog.New(NewFileHandler(ctx, NoopCloser{os.Stderr})), nil
	}

	abspath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%s': %w", fpath, err)
	}

	fd, err := os.OpenFile(abspath, fileFlags, filePerms)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %w", abspath, err)
	}

	return slog.New(NewFileHandler(ctx, fd)), nil
}
