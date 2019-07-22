package accesslog

import "net/http"

type (
	// Entry access log entry
	Entry struct {
		Request  Request
		Response Response
	}

	// Request entry
	Request struct {
		Header http.Header
		Method string
		URL    string
		Body   string
	}

	// Response entry
	Response struct {
		Headers    http.Header
		StatusCode int
		Body       string
	}
)
