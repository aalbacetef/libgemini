# Introduction 

libgemini is a simple Go library that implements support for sending Gemini requests and receiving responses.

It offers a simple API.

Much left to do! 

# geminirc 

The geminirc is meant to be an analogue of the curlrc, but supported by libgemini. 

This enables users to manage settings for apps using libgemini. 
It also makes development a bit simpler.

These settings can also be controlled with environment variables.

The default location for the geminirc file, in order, is:

1. The location specified by the `LIBGEMINI_RC` environment variable.
2. $HOME/.config/libgemini/geminirc 

If it is not found, the directory will be created and the file will be created.

To see a full example check [data/geminirc](data/geminirc)

### Sample geminirc file.

Note: for boolean fields, uncomment to enable, comment to disable.
Comments are set using '#'.


##### Redirects 

Enable this option to automatically follow redirects.

```bash
--follow 
```

##### Trace 

Dump the trace to a file.
Set this option to the path of the file.

Note: that the file will be overwritten on each request.

```bash
--trace /tmp/libgemini-trace.txt
```

##### Dump headers 

 Dump headers of last request to a file.

 Note: that the file will be overwritten on each request.

```bash
--dump-headers /tmp/libgemini-headers.txt
```

##### Insecure mode  

Skip TOFU verification. Overrides --store.

```bash
 --insecure
```

##### Store location 

 Set Store location.

```bash
 --store ~/.config/libgemini/known_hosts
```


 To use an in-memory store:


```bash 
 --store :memory:
```

