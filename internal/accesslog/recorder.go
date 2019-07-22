package accesslog

import (
	"bytes"
	"net/http"
	"router/pkg/net"
)

type responseRecorder struct {
	responseStatusCode    int
	responseBody          bytes.Buffer
	hasBinaryResponseBody bool
	ignoreBody            bool

	http.ResponseWriter
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	size, err := r.ResponseWriter.Write(body)
	if r.hasBinaryResponseBody || err != nil {
		return size, err
	}

	if net.HasBinaryContent(r.ResponseWriter.Header(), body) {
		r.hasBinaryResponseBody = true
		return size, nil
	}

	if !r.ignoreBody {
		r.responseBody.Write(body)
	}

	return size, nil
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.responseStatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
