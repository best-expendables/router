package accesslog

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessLogRecorder(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))
	request.Header.Set("HEADER", "HEADER")

	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.Write([]byte("Response body"))
	recorder.WriteHeader(http.StatusOK)

	request.Header.Set("WRONG-HEADER", "WRONG-HEADER")

	log := recorder.Entry()
	t.Log(log)

	if header := log.Request.Header.Get("WRONG-HEADER"); header != "" {
		t.Error("Headers has been changed")
	}

	if log.Request.Body != "Request body" {
		t.Error("Request body not equals")
	}
	if log.Request.URL != "http://localhost" {
		t.Error("Request URL not eqauls")
	}
	if log.Request.Method != http.MethodGet {
		t.Error("Request method not equals")
	}
	if header := log.Request.Header.Get("HEADER"); header != "HEADER" {
		t.Error("Request headers not equals")
	}
	if log.Response.StatusCode != http.StatusOK {
		t.Error("Response status code not equals")
	}
	if log.Response.Body != "Response body" {
		t.Error("Response body not equals")
	}
	if contentType := log.Response.Headers.Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Error("Response content type not equals")
	}
}

func TestAccessLogRecorer_IgnoreBody(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))

	t.Run("Ignore request body", func(t *testing.T) {
		recorder := New(AccessWriterOptions{
			IgnoreRequestBody: true,
		}, httptest.NewRecorder(), request)
		recorder.Write([]byte("Response body"))

		log := recorder.Entry()
		if log.Request.Body != "" {
			t.Error("Body not empty")
		}
		if log.Response.Body == "" {
			t.Error("Body is an empty")
		}
	})

	t.Run("Ignore response body", func(t *testing.T) {
		recorder := New(AccessWriterOptions{
			IgnoreResponseBody: true,
		}, httptest.NewRecorder(), request)
		recorder.Write([]byte("Response body"))

		log := recorder.Entry()
		if log.Request.Body == "" {
			t.Error("Body is an empty")
		}
		if log.Response.Body != "" {
			t.Error("Body not empty")
		}
	})

	t.Run("Ignore all", func(t *testing.T) {
		recorder := New(AccessWriterOptions{
			IgnoreRequestBody:  true,
			IgnoreResponseBody: true,
		}, httptest.NewRecorder(), request)
		recorder.Write([]byte("Response body"))

		log := recorder.Entry()
		if log.Request.Body != "" {
			t.Error()
		}
		if log.Response.Body != "" {
			t.Error()
		}
	})
}

func TestAccessLogRecorder_RequestBodyIsNil(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	request.Header.Set("HEADER", "HEADER")

	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.Write([]byte("Response body"))
	recorder.WriteHeader(http.StatusOK)

	log := recorder.Entry()
	if log.Request.Body != "" {
		t.Error("Request body should be empty")
	}
}

func TestAccessLogRecorder_ResponseBodyIsNil(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	request.Header.Set("HEADER", "HEADER")

	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.WriteHeader(http.StatusOK)

	log := recorder.Entry()
	if log.Response.Body != "" {
		t.Error("Request body should be empty")
	}
}

// Request and response has content-type for file transfer
func TestAccessLogRecorder_BinaryContentType(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))
	request.Header.Set("Content-Type", "multipart/form-data")

	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.ResponseWriter.Header().Set("Content-Disposition", "attachment")
	recorder.Write([]byte("Response body"))
	recorder.WriteHeader(http.StatusOK)

	log := recorder.Entry()

	if log.Request.Body != "" {
		t.Error("Request body should be empty")
	}

	if log.Response.StatusCode != http.StatusOK {
		t.Error("Response status code not equals")
	}
	if log.Response.Body != "" {
		t.Error("Response body should be empty")
	}
}

// Trying to avoid content type check and send binary with wrong headers
func TestAccessLogRecorder_NoContentType_BinaryRequest(t *testing.T) {
	requestBody := &bytes.Buffer{}

	for _, char := range []rune("Hello world") {
		if err := binary.Write(requestBody, binary.LittleEndian, char); err != nil {
			t.Fatal(err)
		}
	}

	request, _ := http.NewRequest(http.MethodGet, "http://localhost", requestBody)

	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.Write([]byte("Response body"))
	recorder.WriteHeader(http.StatusOK)

	log := recorder.Entry()
	if log.Request.Body != "" {
		t.Error("Request body should be empty!")
	}
}

func TestAccessLogRecorder_NoContentType_Response(t *testing.T) {
	responseBody := &bytes.Buffer{}

	for _, char := range []rune("Hello world") {
		if err := binary.Write(responseBody, binary.LittleEndian, char); err != nil {
			t.Fatal(err)
		}
	}

	request, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	recorder := New(AccessWriterOptions{}, httptest.NewRecorder(), request)
	recorder.Write(responseBody.Bytes())
	log := recorder.Entry()
	if log.Response.Body != "" {
		t.Error("Response body should be empty!")
	}
}
