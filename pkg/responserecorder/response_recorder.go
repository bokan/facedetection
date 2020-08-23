package responserecorder

import (
	"bytes"
	"net/http"
)

type ResponseRecorder struct {
	w          http.ResponseWriter
	r          *http.Request
	buf        *bytes.Buffer
	statusCode int
}

// NewResponseRecorder returns a new instance of write-through response recorder.
func NewResponseRecorder(w http.ResponseWriter, r *http.Request) *ResponseRecorder {
	return &ResponseRecorder{w: w, r: r}
}

// Body returns body content of recorded response.
func (c *ResponseRecorder) Body() []byte {
	return c.buf.Bytes()
}

// StatusCode returns status code of recorded response.
func (c *ResponseRecorder) StatusCode() int {
	return c.statusCode
}

// Header returns header of recorded response.
func (c *ResponseRecorder) Header() http.Header {
	return c.w.Header()
}

// Write should be called inside the HTTP handler.
func (c *ResponseRecorder) Write(bytes []byte) (int, error) {
	if c.buf == nil {
		panic("Write() called before WriteHeader()")
	}
	_, _ = c.buf.Write(bytes)
	return c.w.Write(bytes)
}

// WriteHeader should be called inside the HTTP handler.
func (c *ResponseRecorder) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.w.WriteHeader(statusCode)
	c.buf = bytes.NewBuffer([]byte{})
}
