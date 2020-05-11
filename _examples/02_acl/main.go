package main

import (
	"github.com/best-expendables/logger"
	"github.com/best-expendables/router"
	"github.com/best-expendables/router/middleware"
	"github.com/best-expendables/user-service-client"
	"github.com/go-chi/chi"
	"net/http"
)

var dummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func main() {
	router, err := router.New(router.Configuration{
		LoggerFactory: logger.NewLoggerFactory(logger.InfoLevel),
		NewRelicApp:   &nrmock.NewrelicApp{},
	})

	if err != nil {
		panic(err)
	}

	router.Get("/public", dummyHandler)

	router.Group(func(r chi.Router) {
		r.Use(middleware.OnlyRoles(userclient.RoleAdmin))
	}).Get("/restricted", dummyHandler)
}
