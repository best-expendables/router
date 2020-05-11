package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/best-expendables/logger"
	"github.com/best-expendables/router"
	"github.com/best-expendables/router/middleware"
)

type (
	ErrorResponse struct {
		Message string
		Rvr     interface{}
	}
)

func main() {
	config := router.Configuration{
		LoggerFactory: logger.NewLoggerFactory(logger.DebugLevel),
		PanicHandler:  PanicHandler(),
	}

	mux, _ := router.New(config)
	mux.Get("/", BrokenHandler)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-LEL-User-ID", "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c")
	r.Header.Set("X-LEL-User-NAME", "admin")
	r.Header.Set("X-LEL-User-EMAIL", "admin@lel.com")
	r.Header.Set("X-LEL-User-Roles", "admin,super")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code, "-", w.Body.String())
}

func PanicHandler() middleware.PanicHandler {
	return middleware.PanicHandlerFunc(func(w http.ResponseWriter, r *http.Request, rvr middleware.PanicRecoveredError) {
		response := ErrorResponse{
			Message: "Internal server error",
			Rvr:     rvr,
		}

		data, _ := json.Marshal(&response)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
	})
}

func BrokenHandler(w http.ResponseWriter, r *http.Request) {
	panic("Not enough power, please replace batteries.")
}
