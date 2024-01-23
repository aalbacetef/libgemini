package main

import (
	"flag"
	"fmt"

	"github.com/aalbacetef/libgemini"
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

	storeOpt := libgemini.WithInMemoryStore()
	if storePath != "" {
		storeOpt = libgemini.WithStore(storePath)
	}

	client, err := libgemini.NewClient(
		storeOpt,
		libgemini.WithInsecure(),
	)
	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	fmt.Printf("\n\n%s\n\n", resp.Header.String())

	fmt.Println("--------------------------------")

	fmt.Println(string(resp.Content))
}
