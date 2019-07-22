package accesslog

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"router/pkg/net"
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
		Header: net.CloneHeader(r.Header),
	}

	if r.Body != nil && !opt.IgnoreRequestBody {
		body, _ := ioutil.ReadAll(r.Body)
		if !net.HasBinaryContent(r.Header, body) {
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
