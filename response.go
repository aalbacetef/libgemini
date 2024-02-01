package libgemini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Response struct {
	Header  Header
	MIME    string
	Content []byte
}

type Header struct {
	Meta   string
	Status StatusCode
}

func (hdr Header) String() string {
	return fmt.Sprintf(
		"Status: %s\nMeta: %s",
		hdr.Status.String(), hdr.Meta,
	)
}

const (
	sizeStatusCode       = 2 // Size in bytes.
	maxMetaSize          = 1024
	spaceSize            = 1
	TerminatorSize       = 2 // Size of \r\n.
	minResponseSize      = sizeStatusCode + spaceSize + TerminatorSize
	maxResponseSize      = minResponseSize + maxMetaSize
	spaceByte            = 0x20
	kb                   = 1024
	responsePreallocSize = 4 * kb
)

func ReadResponse(r io.Reader) (Response, error) {
	resp := Response{}

	for {
		p := make([]byte, responsePreallocSize)

		n, err := r.Read(p)
		if err != nil && !errors.Is(err, io.EOF) {
			return resp, fmt.Errorf("error reading response: %w", err)
		}

		p = p[:n]

		if n == 0 {
			return resp, nil
		}

		if resp.Header.Status == Unset {
			header, bytesRead, parseErr := parseHeader(p)
			if parseErr != nil {
				return resp, fmt.Errorf("error parsing header: %w", parseErr)
			}

			p = p[bytesRead:]

			resp.Header = header
			resp.MIME = mimeFromHeader(header)
		}

		resp.Content = append(resp.Content, p...)

		if errors.Is(err, io.EOF) {
			return resp, nil
		}
	}
}

const DefaultMimeType = "text/gemini; charset=utf-8"

func parseHeader(respBytes []byte) (Header, int, error) {
	index := bytes.Index(respBytes, []byte{'\r', '\n'})
	if index == -1 {
		return Header{}, 0, fmt.Errorf("no CRLF found")
	}

	n := index + TerminatorSize
	p := respBytes[:index]

	hdr := Header{}

	spaceIndex := bytes.IndexByte(p, spaceByte)
	if spaceIndex == -1 {
		return Header{}, 0, fmt.Errorf("could not find space")
	}

	code, err := strconv.Atoi(string(p[:spaceIndex]))
	if err != nil {
		return hdr, 0, fmt.Errorf("could not parse status code: %w", err)
	}

	hdr.Status = StatusCode(code)
	hdr.Meta = string(p[spaceIndex+1:])

	metaBytes := len(hdr.Meta)
	if metaBytes > maxMetaSize {
		return hdr, 0, fmt.Errorf("max meta size of %d bytes exceeded, got %d bytes", maxMetaSize, metaBytes)
	}

	return hdr, n, nil
}

func mimeFromHeader(hdr Header) string {
	if !hdr.Status.IsSuccess() {
		return ""
	}

	m := hdr.Meta
	if m == "" {
		m = DefaultMimeType
	}

	return m
}
