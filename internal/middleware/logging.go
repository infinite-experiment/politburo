// middleware/logging.go
package middleware

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type respLogger struct {
	http.ResponseWriter
	status int
	buf    *bytes.Buffer
}

func (l *respLogger) WriteHeader(code int) {
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *respLogger) Write(b []byte) (int, error) {
	l.buf.Write(b)
	return l.ResponseWriter.Write(b)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("→ %s %s", r.Method, r.URL.String())

		for name, vals := range r.Header {
			for _, v := range vals {
				log.Printf("→ Req Header: %s=%s", name, v)
			}
		}

		// wrap response
		buf := &bytes.Buffer{}
		lw := &respLogger{ResponseWriter: w, status: http.StatusOK, buf: buf}

		log.Printf("→ %s %s", r.Method, r.URL.Path)

		start := time.Now()
		next.ServeHTTP(lw, r)
		dur := time.Since(start)

		log.Printf("← %d %s (%s)\nResponse Body: %s",
			lw.status, http.StatusText(lw.status), dur, buf.String(),
		)
	})
}
