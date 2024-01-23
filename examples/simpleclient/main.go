package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aalbacetef/libgemini"
)

func main() {
	url := ""

	flag.StringVar(&url, "url", url, "URL to request (doesn't need to be prefixed with 'gemini://')")
	flag.Parse()

	if url == "" {
		flag.Usage()
		fmt.Println("error: please provide a URL")

		return
	}

	// Use libgemini.NewClient to pass in options.
	// A few helper functions are provided.
	client, err := libgemini.NewClient(
		func(opts *libgemini.Options) {
			opts.DumpHeaders = "./dump-headers.json.log"
			opts.Trace = "./trace.json.log"
			opts.Timeout = 1 * time.Minute
		},
		libgemini.WithInMemoryStore(),
	)
	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	// The API should be familiar to users of net/http.
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	if !resp.Header.Status.IsSuccess() {
		fmt.Println("request failed")
		fmt.Println("hint: pipe the *.log files into jq or print resp.Header")
	}

	fmt.Println(string(resp.Content))
}
