package accesslog

import (
	"bufio"
	"bytes"
	internalNet "github.com/best-expendables/router/pkg/net"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
)

type (
	AccessWriterOptions struct {
		IgnoreRequestBody  bool
		IgnoreResponseBody bool
	}

	AccessWriter struct {
		responseRecorder
		request Request
	}
)

func New(opt AccessWriterOptions, w http.ResponseWriter, r *http.Request) *AccessWriter {
	requestEntry := Request{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: internalNet.CloneHeader(r.Header),
	}

	if r.Body != nil && !opt.IgnoreRequestBody {
		body, _ := ioutil.ReadAll(r.Body)
		if !internalNet.HasBinaryContent(r.Header, body) {
			requestEntry.Body = string(body)
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	return &AccessWriter{
		request: requestEntry,
		responseRecorder: responseRecorder{
			ResponseWriter: w,
			ignoreBody:     opt.IgnoreResponseBody,
		},
	}
}

// Entry returns access log entry
func (w *AccessWriter) Entry() Entry {
	return Entry{
		Request: w.request,
		Response: Response{
			Headers:    w.Header(),
			StatusCode: w.responseStatusCode,
			Body:       w.responseBody.String(),
		},
	}
}

func (w *AccessWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.responseRecorder.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("chi/middleware: http.Hijacker is unavailable on the writer")
}
