package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"io/ioutil"

	"bitbucket.org/snapmartinc/logger"
	"bitbucket.org/snapmartinc/trace"
)

type LogContent struct {
	TraceID   string `json:"traceId"`
	UserID    string `json:"userId"`
	ContextID string `json:"context-id"`
	Content   struct {
		Request  map[string]string `json:"request"`
		Response map[string]string `json:"response"`
	} `json:"content"`
}

func TestAccessLog(t *testing.T) {
	buf := &bytes.Buffer{}

	entry := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf)).
		Logger(context.TODO())

	requestBody := "REQUEST BODY"
	request, _ := http.NewRequest(http.MethodGet, "http://localhost",
		bytes.NewBufferString(requestBody))
	request.Header.Set("HEADER_TEST", "HEADER_TEST")
	trace.RequestIDToHeader(request.Header, "REQUEST-ID")

	request = request.WithContext(logger.ContextWithEntry(entry, request.Context()))

	responseBody := "RESPONSE BODY"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ioutil.ReadAll(r.Body)

		// This header should be ignore into the access log
		r.Header.Set("INVALID_HEADER", "INVALID_HEADER")

		w.Write([]byte(responseBody))
		w.WriteHeader(http.StatusOK)
	})

	fnext := AccessLog(AccessLogOptions{})
	fnext(next).ServeHTTP(httptest.NewRecorder(), request)

	if buf.Len() == 0 {
		t.Fatal("Buffer is empty")
	}

	log := buf.String()
	t.Log(log)

	if !strings.Contains(log, requestBody) {
		t.Error("Log does not contains request body")
	}
	if !strings.Contains(log, responseBody) {
		t.Error("Log does not contains response body")
	}
	if !strings.Contains(log, "\"StatusCode\":200") {
		t.Error("Log does not contains status code")
	}
	if !strings.Contains(log, "REQUEST-ID") {
		t.Error("Log does not contains request id")
	}
	if strings.Contains(log, "INVALID_HEADER") {
		t.Error("Headers mutable")
	}
}

func TestAccessLog_StatusUnprocessableEntity(t *testing.T) {
	buf := &bytes.Buffer{}

	entry := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf)).
		Logger(context.TODO())

	request, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	request = request.WithContext(logger.ContextWithEntry(entry, request.Context()))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	fnext := AccessLog(AccessLogOptions{})
	fnext(next).ServeHTTP(httptest.NewRecorder(), request)

	log := buf.String()
	t.Log(log)

	if !strings.Contains(log, "\"level\":\"info\"") {
		t.Error("Should be info level")
	}

	if !strings.Contains(log, http.StatusText(http.StatusUnprocessableEntity)) {
		t.Error()
	}
}

func TestAccessLog_IgnoreBody(t *testing.T) {
	t.Run("AccessLog ignore request body", func(t *testing.T) {
		buf := &bytes.Buffer{}
		entry := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf)).
			Logger(context.TODO())
		request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))
		request = request.WithContext(logger.ContextWithEntry(entry, request.Context()))
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ioutil.ReadAll(r.Body)
			w.Write([]byte("Response body"))
			w.WriteHeader(http.StatusOK)
		})
		fnext := AccessLog(AccessLogOptions{IgnoreRequestBody: true})
		fnext(next).ServeHTTP(httptest.NewRecorder(), request)

		log := buf.String()
		var logContent LogContent
		json.Unmarshal([]byte(log), &logContent)

		if logContent.Content.Request["Body"] != "" {
			t.Error("Request body should be ignored")
		}
		if logContent.Content.Response["Body"] == "" {
			t.Error("Response body should not be ignored")
		}
	})

	t.Run("AccessLog ignore response body", func(t *testing.T) {
		buf := &bytes.Buffer{}
		entry := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf)).
			Logger(context.TODO())
		request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))
		request = request.WithContext(logger.ContextWithEntry(entry, request.Context()))
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ioutil.ReadAll(r.Body)
			w.Write([]byte("Response body"))
			w.WriteHeader(http.StatusOK)
		})
		fnext := AccessLog(AccessLogOptions{IgnoreResponseBody: true})
		fnext(next).ServeHTTP(httptest.NewRecorder(), request)

		log := buf.String()
		var logContent LogContent
		json.Unmarshal([]byte(log), &logContent)

		if logContent.Content.Request["Body"] == "" {
			t.Error("Request body should not be ignored")
		}
		if logContent.Content.Response["Body"] != "" {
			t.Error("Response body should be ignored")
		}

		t.Log(log)
	})

	t.Run("AccessLog ignore all", func(t *testing.T) {
		buf := &bytes.Buffer{}
		entry := logger.NewLoggerFactory(logger.DebugLevel, logger.SetOut(buf)).
			Logger(context.TODO())
		request, _ := http.NewRequest(http.MethodGet, "http://localhost", bytes.NewBufferString("Request body"))
		request = request.WithContext(logger.ContextWithEntry(entry, request.Context()))
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ioutil.ReadAll(r.Body)
			w.Write([]byte("Response body"))
			w.WriteHeader(http.StatusOK)
		})
		fnext := AccessLog(AccessLogOptions{IgnoreRequestBody: true, IgnoreResponseBody: true})
		fnext(next).ServeHTTP(httptest.NewRecorder(), request)

		log := buf.String()
		var logContent LogContent
		json.Unmarshal([]byte(log), &logContent)

		if logContent.Content.Request["Body"] != "" {
			t.Error("Request body should be ignored")
		}
		if logContent.Content.Response["Body"] != "" {
			t.Error("Response body should be ignored")
		}

		t.Log(log)
	})
}
