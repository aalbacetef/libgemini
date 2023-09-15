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
	Content []byte
	MIME    string
}

// HEADER

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
	sizeStatusCode       = 2 // in bytes
	maxMetaSize          = 1024
	spaceSize            = 1
	EOFSize              = 2 // size of \r\n
	minResponseSize      = sizeStatusCode + spaceSize + EOFSize
	maxResponseSize      = minResponseSize + maxMetaSize
	spaceByte            = 0x20
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
		log.Println("len(p): ", len(p))

		p = p[:n]
		log.Println("n: ", n)

		if n == 0 {
			return resp, nil
		}

		if !parsedHeader {
			index := bytes.Index(p, []byte{'\r', '\n'})
			if index == -1 {
				return resp, fmt.Errorf("no CRLF found")
			}

			headerBytes := p[:index+2]
			header, parseErr := parseHeader(headerBytes)
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
			log.Println("io.EOF")
			return resp, nil
		}
	}
}

func parseHeader(p []byte) (Header, error) {
	h := Header{}
	spaceIndex := bytes.IndexByte(p, spaceByte)
	if spaceIndex == -1 {
		return Header{}, fmt.Errorf("could not find space")
	}

	code, err := strconv.Atoi(string(p[:spaceIndex]))
	if err != nil {
		return h, fmt.Errorf("could not parse status code: %w", err)
	}

	h.Status = StatusCode(code)
	h.Meta = string(p[spaceIndex+1:])
	metaBytes := len(h.Meta)
	if metaBytes > maxMetaSize {
		return h, fmt.Errorf("max meta size of %d bytes exceeded, got %d bytes", maxMetaSize, metaBytes)
	}

	return h, nil
}

/*
// ReadHeader will first read the status code and space, run checks on that,
// then proceed to check if there is a meta to read.
// It will ensure that the following validation rules are followed:
//   - the header doesn't start with BOM (0xFEFF)
//   - the meta field doesn't start with (0xFEFF)
//   - the status code is followed by a space
//   - the response ends with CRLF (\r\n)
func ReadHeader(r io.Reader) (Header, error) {
	p := make([]byte, sizeStatusCode+spaceSize)
	h := Header{}

	// read status code and space
	n, err := r.Read(p)
	if n != len(p) {
		return h, fmt.Errorf("insufficient bytes read: %d, expected %d", n, len(p))
	}
	if err != nil {
		return h, fmt.Errorf("")
	}

	// rule: header must not start with BOM
	if err := mustNotStartHeaderWithBOM(p); err != nil {
		return h, err
	}

	// parse code into uint
	code, err := strconv.ParseUint(string(p[:sizeStatusCode]), 10, 64)
	if err != nil {
		return h, fmt.Errorf("could not parse status code: %w", err)
	}

	h.Status = StatusCode(code)

	// rule: status code must be followed by a space
	if p[n-1] != spaceByte {
		return h, fmt.Errorf("status code not followed by a space, got %#0x", p[n-1])
	}

	// check if there is a meta or a CRLF

}*/

// ReadHeader will allocate enough space to read a response of the maximum size defined by the protocol.
// It will apply the following validation on the request:
//   - ensure the header doesn't start with 0xFEFF
//   - ensure the META field doesn't start with 0xFEFF
//   - ensure the status code is followed by a space
//   - ensure the request ends with CRLF (\r\n)
//
// see "3 Gemini responses" of the gemini specification.

/*
func ReadHeader(r io.Reader) (Header, error) {
	p := make([]byte, maxResponseSize)

	n, err := r.Read(p)
	if err != nil {
		return Header{}, fmt.Errorf("could not read response: %w", err)
	}

	for _, rule := range []ruleFunc{
		mustHaveMinRequestSize,
		mustNotStartHeaderWithBOM,
	} {
		if err := rule(p); err != nil {
			return Header{}, err
		}
	}

	code, err := strconv.ParseUint(string(p[:sizeStatusCode]), 10, 64)
	if err != nil {
		return Header{}, fmt.Errorf("could not parse status code: %w", err)
	}

	h := Header{
		Status: StatusCode(code),
		Meta:   string(p[spaceIndex+1 : n-2]),
	}

	for _, rule := range []ruleFunc{
		mustHaveSpaceAfterStatus,
		mustEndWithCRLF,
	} {
		if err := rule(p); err != nil {
			return h, err
		}
	}

	return h, nil
}

type ruleFunc func([]byte) error

// mustHaveMinRequestSize checks that the request has the minimum required size.
func mustHaveMinRequestSize(p []byte) error {
	n := len(p)
	if n < minRequestSize {
		return fmt.Errorf("insufficient bytes, read %d minimum: %d", n, minRequestSize)
	}

	return nil
}

const (
	BOMLow  = 0xFE
	BOMHigh = 0xFF
	BOMLen  = 2
)

func isBOM(p []byte) bool {
	if len(p) != BOMLen {
		return false
	}

	return p[0] == BOMLow && p[1] == BOMHigh
}

// mustNotStartHeaderWithBOM checks that the header doesn't start with the BOM (0xFEFF).
func mustNotStartHeaderWithBOM(p []byte) error {
	const size = BOMLen
	n := len(p)

	if n < size {
		return fmt.Errorf("header too small, expected at least %d bytes, got %d", size, n)
	}

	if isBOM(p[:size]) {
		return fmt.Errorf("meta field MUST not start with BOM")
	}

	return nil
}

// mustHaveSpaceAfterStatus checks that the space byte (0x20) occurs after a status code.
func mustHaveSpaceAfterStatus(p []byte) error {
	if p[spaceIndex] != spaceByte {
		return fmt.Errorf("expected space(%#0x), got %#0x", spaceByte, p[spaceIndex])
	}

	return nil
}

// mustEndWithCRLF validates that the response ends with \r\n.
func mustEndWithCRLF(p []byte) error {
	const size = 2
	n := len(p)
	tail := p[n-2 : n]

	fmt.Printf("%v\n", p)
	if tail[0] != '\r' || tail[1] != '\n' {
		return fmt.Errorf("request MUST end with CRLF")
	}

	return nil
}*/

// BODY

const (
	kb               = 1024
	bodyPreallocSize = 4 * kb
)

/*
// ReadBody will read a response body. It preallocates 4Kb and reads the
// response in chunks of 1Kb.
func ReadBody(r io.Reader) ([]byte, error) {
	read := make([]byte, 0, bodyPreallocSize)

	for {
		data := make([]byte, kb)

		n, err := r.Read(data)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("got error while reading: %w", err)
		}

		if n == 0 {
			return read, nil
		}

		p := data[:n]
		read = append(read, p...)
		if err == io.EOF {
			return read, nil
		}
	}
}*/
