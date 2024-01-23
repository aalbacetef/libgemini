package libgemini

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type FileHandler struct {
	fd io.WriteCloser
	slog.Handler
}

func NewFileHandler(ctx context.Context, fd io.WriteCloser) FileHandler {
	go func() {
		<-ctx.Done()
		fd.Close()
	}()

	return FileHandler{
		fd:      fd,
		Handler: slog.NewTextHandler(fd, nil),
	}
}

type NoopHandler struct{}

func (NoopHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (NoopHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (NoopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NoopHandler{}
}

func (NoopHandler) WithGroup(string) slog.Handler {
	return NoopHandler{}
}

const (
	fileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	filePerms = fs.FileMode(0o644)
)

// NewLoggerFromPath will take in a path and if it's empty it will
// return a logger using a NoopHandler, otherwise it will try to
// initialize a FileHandler based on the filepath. it is
// expected to be used alongside a context.
func NewLoggerFromPath(ctx context.Context, fpath string) (*slog.Logger, error) {
	if fpath == "" {
		return slog.New(NoopHandler{}), nil
	}

	abspath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, err
	}

	fd, err := os.OpenFile(abspath, fileFlags, filePerms)
	if err != nil {
		return nil, err
	}

	return slog.New(NewFileHandler(ctx, fd)), nil
}
