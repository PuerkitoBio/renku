package cache

import (
	"bytes"
	"io"
	"net/http"
)

type cacheWriter struct {
	http.ResponseWriter
	code  int
	b     *bytes.Buffer
	multi io.Writer
}

func (ø *cacheWriter) WriteHeader(code int) {
	ø.code = code
	ø.ResponseWriter.WriteHeader(code)
}

func (ø *cacheWriter) Write(b []byte) (int, error) {
	return ø.multi.Write(b)
}

func LRUCacheHandler(h http.Handler) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO : Check if available in cache, if so, done.

		// Create the cache writer
		buf := bytes.NewBuffer(nil)
		cw := &cacheWriter{
			w,
			http.StatusOK,
			buf, // May not be required in the cache struct
			io.MultiWriter(w, buf),
		}

		// Call the wrapped handler with the cache writer
		h.ServeHTTP(cw, r)

		if cw.code >= 200 && cw.code < 300 {
			// TODO : Store the response in the cache
			// TODO : Save the header's keys too
		}
	}
}
