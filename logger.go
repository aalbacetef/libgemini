package libgemini

import (
	"context"
	"io"
	"log/slog"
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
