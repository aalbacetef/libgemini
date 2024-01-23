package libgemini

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"
)

type mockCloser struct {
	io.ReadWriter
	closeCh chan<- struct{}
}

func (mc mockCloser) Close() error {
	mc.closeCh <- struct{}{}

	return nil
}

const (
	testTimeout = 5 * time.Second
	testMsg     = "test message"
)

func TestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	buf := &bytes.Buffer{}
	ch := make(chan struct{})

	h := NewFileHandler(ctx, mockCloser{ReadWriter: buf, closeCh: ch})
	l := slog.New(h)

	t.Run("it writes to the logger", func(tt *testing.T) {
		l.Info(testMsg, "buf", buf)
		if !strings.Contains(buf.String(), testMsg) {
			tt.Fatalf("did not find test message '%s' in buffer", testMsg)
		}
	})

	cancel()

	t.Run("it calls Close on context canceled", func(tt *testing.T) {
		select {
		case <-time.After(testTimeout):
			tt.Fatalf("test timed out waiting %s", testTimeout)
		case <-ch:
			return
		}
	})
}
