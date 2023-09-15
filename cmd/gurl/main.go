package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/aalbacetef/libgemini"
)

func main() {
	out := ""
	headersOnly := false

	flag.BoolVar(&headersOnly, "headers-only", headersOnly, "only print headers")
	flag.StringVar(&out, "out", out, "output file path")

	flag.Parse()

	nargs := flag.NArg()
	if nargs == 0 {
		flag.Usage()
		return
	}

	w := os.Stdout
	if out != "" {
		fpath := filepath.Clean(out)
		fd, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			fmt.Printf("could not open file (%s): %v\n", fpath, out)
			return
		}
		w = fd
		defer fd.Close()
	}

	args := flag.Args()
	URL := args[0]

	client := libgemini.Client{
		Timeout: 30 * time.Second,
		Config: &tls.Config{
			MinVersion:         tls.VersionTLS13,
			InsecureSkipVerify: true,
		},
	}

	resp, err := client.Get(URL)
	if err != nil {
		fmt.Println("error getting page: ", err)
		return
	}

	fmt.Fprintln(w, resp.Header.String())
	fmt.Fprintln(w, string(resp.Content))
}
