package libgemini

import (
	_ "embed"
	"testing"
)

//go:embed testdata/geminirc
var rcTestFile []byte

func TestConfigFile(t *testing.T) {
	opts := configOpts(string(rcTestFile))
	want := map[string]strOrBool{
		KeyFollowRedirects: {b: true},
		KeyTrace:           {s: "/tmp/libgemini-trace.txt"},
		KeyInsecure:        {b: true},
		KeyStorePath:       {s: "~/.config/libgemini/known_hosts"},
	}

	for key, val := range want {
		got, ok := opts[key]
		if !ok {
			t.Fatalf("(%s) not set", key)
		}

		if got.s != val.s {
			t.Fatalf("(%s) got %s, want %s", key, got.s, val.s)
		}

		if got.b != val.b {
			t.Fatalf("(%s) got %t, want %t", key, got.b, val.b)
		}
	}
}
