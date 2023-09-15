package libgemini

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"
)

const (
	testPort    = 1965
	testTimeout = 5 * time.Second
)

func TestReadResponse(t *testing.T) {
	const (
		testHost = "gemini.circumlunar.space"
		testPath = "/docs/specification.gmi"
		testURL  = testHost + testPath
	)

	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		ServerName:         testHost,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", testHost, testPort), cfg)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer conn.Close()

	t.Logf("connection established")

	c := Client{
		Config:  cfg,
		Timeout: testTimeout,
	}

	resp, err := c.Get(testURL)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Logf("response: \n%+v", resp.Header)
	t.Logf("content length: %d", len(resp.Content))
}

func TestResponse2(t *testing.T) {
	const (
		testHost = "geminiprotocol.net"
		testPath = "/docs/specification.gmi"
		testURL  = testHost + testPath
	)

	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		ServerName:         testHost,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", testHost, testPort), cfg)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer conn.Close()

	t.Logf("connection established")

	c := Client{
		Config:  cfg,
		Timeout: testTimeout,
	}

	resp, err := c.Get(testURL)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Logf("response: \n%+v", resp.Header)
	t.Logf("content length: %d", len(resp.Content))
}
