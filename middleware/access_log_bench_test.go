package middleware

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/best-expendables/logger"
)

func BenchmarkAccessLog(b *testing.B) {
	e := logger.NewLoggerFactory(
		logger.InfoLevel,
		logger.SetOut(ioutil.Discard),
	).Logger(context.TODO())

	req, _ := http.NewRequest(http.MethodGet, "http://localhost/v1/user/f77ea2d3-d040-4e05-bd43-48893743b2a9", nil)
	req = req.WithContext(
		logger.ContextWithEntry(e, context.Background()),
	)

	req.Header.Add("Header", "Header")
	req.Header.Add("Role", "admin")
	req.Header.Add("Trace-ID", "f77ea2d3-d040-4e05-bd43-48893743b2a9")
	req.Header.Add("Context-ID", "c3a433e9-d5a8-44eb-b97d-a9f773680201")

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello world!"))
		w.WriteHeader(http.StatusOK)
	})

	fnext := AccessLog(AccessLogOptions{})
	h2 := fnext(h)

	b.Run("Concurrent", func(b *testing.B) {
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				h2.ServeHTTP(httptest.NewRecorder(), req)
			}
		})
	})
}
