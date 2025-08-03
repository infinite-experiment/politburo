package common

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func LogHTTPRequest(req *http.Request) {
	// Make a copy of the body if it exists
	var bodyCopy []byte
	if req.Body != nil {
		bodyCopy, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // reset body
	}

	dump, err := httputil.DumpRequestOut(req, true) // true to include body
	if err != nil {
		log.Printf("Failed to dump HTTP request: %v", err)
	} else {
		log.Printf("HTTP Request Dump:\n%s", dump)
	}

	// Reset the body again (req.Body may be read again later)
	if bodyCopy != nil {
		req.Body = io.NopCloser(bytes.NewReader(bodyCopy))
	}
}
