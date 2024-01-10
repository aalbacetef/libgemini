package main

import (
	"flag"
	"fmt"

	"github.com/aalbacetef/libgemini"
	"github.com/aalbacetef/libgemini/tofu"
)

func main() {
	url := ""
	storePath := ""

	flag.StringVar(&url, "url", url, "URL to request (doesn't need to be prefixed with 'gemini://')")
	flag.StringVar(&storePath, "store", storePath, "path to known_hosts file")
	flag.Parse()

	if url == "" {
		flag.Usage()
		fmt.Println("error: please provide a URL")
		return
	}

	var store tofu.Store = tofu.NewInMemoryStore()

	if storePath != "" {
		_store, err := tofu.NewFileStore(storePath)
		if err != nil {
			fmt.Println("error: ", err)

			return
		}

		store = _store
	}

	client := libgemini.NewClient(store)

	// check a URL
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	// the header
	fmt.Printf("\n\n%s\n\n", resp.Header.String())

	// divisor
	fmt.Println("--------------------------------")

	// print the body
	fmt.Println(string(resp.Content))
}
