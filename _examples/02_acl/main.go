package main

import (
	"bitbucket.org/snapmartinc/logger"
	"bitbucket.org/snapmartinc/router"
	"bitbucket.org/snapmartinc/router/middleware"
	"bitbucket.org/snapmartinc/user-service-client"
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
