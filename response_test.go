package libgemini

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type testCase = struct {
	label   string
	rawResp []byte
	want    []byte
}

func TestReadResponse(t *testing.T) {
	cases := []testCase{
		{
			"response with MIME",
			mustReadFile(t, "testdata/response.raw"),
			mustReadFile(t, "testdata/response.json"),
		},
		{
			"response no MIME",
			mustReadFile(t, "testdata/response.no-mime.raw"),
			mustReadFile(t, "testdata/response.no-mime.json"),
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(tt *testing.T) {
			runTest(tt, c)
		})
	}
}

func runTest(t *testing.T, c testCase) {
	t.Helper()

	r := bytes.NewReader(c.rawResp)

	resp, err := ReadResponse(r)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	want := trim(c.want)

	got, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatalf("could not encode JSON: %v", err)
	}

	wantR := Response{}

	if err := json.Unmarshal(want, &wantR); err != nil {
		t.Fatalf("could not unmarshal json: %v", err)
	}

	m, n := len(got), len(want)
	if m != n {
		t.Fatalf("(length) got %d bytes, want %d", m, n)
	}

	for k, b := range got {
		if b != want[k] {
			t.Fatalf(
				"responses differ. got.Content: '%s', want: '%s'",
				string(resp.Content), string(wantR.Content),
			)
		}
	}
}

func trim(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func mustReadFile(t *testing.T, fpath string) []byte {
	t.Helper()

	data, err := os.ReadFile(fpath)
	if err != nil {
		t.Fatalf("could not read '%s': %v", fpath, err)
	}

	return data
}
