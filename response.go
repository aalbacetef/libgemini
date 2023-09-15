package libgemini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Response struct {
	Header  Header
	MIME    string
	Content []byte
}

type Header struct {
	Status StatusCode
	Meta   string
}

func (hdr Header) String() string {
	return fmt.Sprintf(
		"\nStatus: %s (is successful: %t)\nMeta: %s\n",
		hdr.Status.String(), hdr.Status.IsSuccess(), hdr.Meta,
	)
}

const (
	sizeStatusCode       = 2 // Size in bytes.
	maxMetaSize          = 1024
	spaceSize            = 1
	EOFSize              = 2 // Size of \r\n.
	minResponseSize      = sizeStatusCode + spaceSize + EOFSize
	maxResponseSize      = minResponseSize + maxMetaSize
	spaceByte            = 0x20
	kb                   = 1024
	responsePreallocSize = 4 * kb
)

func ReadResponse(r io.Reader) (Response, error) {
	bytesRead := 0
	parsedHeader := false
	resp := Response{}

	for {
		p := make([]byte, responsePreallocSize)

		n, err := r.Read(p)
		if err != nil && !errors.Is(err, io.EOF) {
			return resp, fmt.Errorf("error reading response: %w", err)
		}

		bytesRead += n
		p = p[:n]

		if n == 0 {
			return resp, nil
		}

		if !parsedHeader {
			index := bytes.Index(p, []byte{'\r', '\n'})
			if index == -1 {
				return resp, fmt.Errorf("no CRLF found")
			}

			header, parseErr := parseHeader(p[:index+2])
			if parseErr != nil {
				return resp, fmt.Errorf("error parsing header: %w", parseErr)
			}

			resp.Header = header
			parsedHeader = true
			p = p[index+2:]
		}

		log.Println("read bytes: ", bytesRead)
		log.Println("p: ", string(p))
		resp.Content = append(resp.Content, p...)

		if errors.Is(err, io.EOF) {
			return resp, nil
		}
	}
}

func parseHeader(p []byte) (Header, error) {
	hdr := Header{}

	spaceIndex := bytes.IndexByte(p, spaceByte)
	if spaceIndex == -1 {
		return Header{}, fmt.Errorf("could not find space")
	}

	code, err := strconv.Atoi(string(p[:spaceIndex]))
	if err != nil {
		return hdr, fmt.Errorf("could not parse status code: %w", err)
	}

	hdr.Status = StatusCode(code)
	hdr.Meta = string(p[spaceIndex+1:])

	metaBytes := len(hdr.Meta)
	if metaBytes > maxMetaSize {
		return hdr, fmt.Errorf("max meta size of %d bytes exceeded, got %d bytes", maxMetaSize, metaBytes)
	}

	return hdr, nil
}
