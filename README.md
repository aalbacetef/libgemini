![CI status](https://github.com/aalbacetef/libgemini/actions/workflows/ci.yml/badge.svg) [![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

# Libgemini

## Introduction 

Libgemini is a simple Gemini client library for Go, allowing you to interact with Gemini servers and retrieve content over the Gemini protocol.

Support for writing servers and gemtext parsing is on the roadmap.


#### Features

- Supports a geminirc file as well as environment variables for controlling behavior.
- Simple and easy-to-use API for interacting with Gemini servers.
- 0 dependencies, only stdlib.

## Installation 

```bash
$ go get -u github.com/aalbacetef/libgemini
```

## Usage 

For a working example, see [examples/simpleclient](examples/simpleclient).


## geminirc 

The geminirc is meant to be an analogue of the curlrc, but supported by libgemini. 

This enables users to manage settings for apps using libgemini. 
It also makes development a bit simpler.

These settings can also be controlled with environment variables.

The default location for the geminirc file, in order, is:

1. The location specified by the `LIBGEMINI_RC` environment variable.
2. $HOME/.config/libgemini/geminirc 

If it is not found, the directory will be created and the file will be created.

To see a full example check [data/geminirc](data/geminirc)


### Environment Variables

The following environment variables are supported:

 - `LIBGEMINI_RC`
 - `LIBGEMINI_FOLLOW_REDIRECTS`
 - `LIBGEMINI_STORE_PATH`
 - `LIBGEMINI_DUMP_HEADERS`
 - `LIBGEMINI_TRACE`
 - `LIBGEMINI_INSECURE`

See the below section for their usage.


### Sample geminirc file.

Note: for boolean fields, uncomment to enable, comment to disable.
Comments are set using '#'.


##### Redirects 

Env: `LIBGEMINI_FOLLOW_REDIRECTS`

Enable this option to automatically follow redirects.

```bash
--follow 
```

##### Trace 

Env: `LIBGEMINI_TRACE`

Dump the trace to a file.
Set this option to the path of the file.

Note: that the file will be overwritten on each request.

```bash
--trace /tmp/libgemini-trace.txt
```

##### Dump headers 

Env: `LIBGEMINI_DUMP_HEADERS`

Dump headers of last request to a file.

Note: that the file will be overwritten on each request.

```bash
--dump-headers /tmp/libgemini-headers.txt
```

##### Insecure mode  

Env: `LIBGEMINI_INSECURE`

Skip TOFU verification. Overrides --store.

```bash
 --insecure
```

##### Store location 

Env: `LIBGEMINI_STORE_PATH`

Set Store location.

```bash
 --store ~/.config/libgemini/known_hosts
```


To use an in-memory store:


```bash 
 --store :memory:
```


## Contributing 

If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

