package httpcache

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

func (c *ResponseRecorder) Body() []byte {
	return c.buf.Bytes()
}

func (c *ResponseRecorder) StatusCode() int {
	return c.statusCode
}

func NewResponseRecorder(w http.ResponseWriter, r *http.Request) *ResponseRecorder {
	return &ResponseRecorder{w: w, r: r}
}

func (c *ResponseRecorder) Header() http.Header {
	return c.w.Header()
}

func (c *ResponseRecorder) Write(bytes []byte) (int, error) {
	if c.buf == nil {
		panic("Write() called before WriteHeader()")
	}
	_, _ = c.buf.Write(bytes)
	return c.w.Write(bytes)
}

func (c *ResponseRecorder) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.w.WriteHeader(statusCode)
	c.buf = bytes.NewBuffer([]byte{})
}
