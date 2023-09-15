package libgemini

import "testing"

func TestRequest(t *testing.T) {
	cases := []struct {
		url   string
		want  string
		label string
	}{
		{
			"geminiprotocol.net",
			"gemini://geminiprotocol.net",
			"no scheme, no path",
		},
		{
			"geminiprotocol.net/docs",
			"gemini://geminiprotocol.net/docs",
			"no scheme, with path",
		},
		{
			"gemini://geminiprotocol.net/docs",
			"gemini://geminiprotocol.net/docs",
			"with scheme, with path",
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(tt *testing.T) {
			req, err := NewRequest(c.url)
			if err != nil {
				tt.Fatalf("error creating request: %v", err)
			}

			got := req.String()
			want := c.want

			if got != want {
				tt.Fatalf("got %s, want %s", got, want)
			}
		})
	}
}
