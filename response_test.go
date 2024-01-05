package libgemini

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"
)

//go:embed testdata/response.raw
var testRawResponse []byte

//go:embed testdata/response.json
var testResponseJSON []byte

func TestReadResponse(t *testing.T) {
	r := bytes.NewReader(testRawResponse)

	resp, err := ReadResponse(r)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	want := testResponseJSON

	got, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatalf("could not encode JSON: %v", err)
	}

	m, n := len(got), len(want)
	if m != n {
		t.Fatalf("(length) got %d bytes, want %d", m, n)
	}

	for k, b := range got {
		if b != want[k] {
			t.Fatalf("responses differ. got.Content: %s", string(resp.Content))
		}
	}
}
