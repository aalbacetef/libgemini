

# Introduction 

This example shows how to use the library to build a client with options.

# Usage 

You can run the client by either directly using `go run`:

```bash
$ go run ./examples/simpleclient/ -url geminiprotocol.net/ 
```

OR building it:

```bash
$ go build ./examples/simpleclient/
$ ./simpleclient -url geminiprotocol.net/
```

# Notes

## Trailing slash 

If you're getting back a `redirect` response, the most likely scenario is you 
forgot to add a trailing slash. You can confirm this by checking the Meta field.

Check this out by passing in `-url geminiprotocol.net`.

