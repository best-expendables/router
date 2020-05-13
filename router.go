package router

import (
	"errors"
	"github.com/best-expendables/router/middleware"
	newrelic "github.com/newrelic/go-agent"
	"net/http"
	"os"

	"github.com/best-expendables/logger"

	"github.com/go-chi/chi"
)

type (
	// Configuration router configuration
	Configuration struct {
		// LoggerFactory using in ContextLogger middleware
		LoggerFactory logger.Factory

		// NewrelicApp
		NewrelicApp newrelic.Application

		// PanicHandler optional parameter
		// On nil panic returns only 500 status code
		PanicHandler middleware.PanicHandler

		AccessLog struct {
			// Disable access log middleware
			Disable bool
		}
	}
)

// New returns router with default list of middlewares
func New(cfg Configuration) (chi.Router, error) {
	if cfg.LoggerFactory == nil {
		return nil, errors.New("logger factory is missing")
	}

	prefix, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	middlewares := make([]func(http.Handler) http.Handler, 0)
	middlewares = append(middlewares, middleware.RequestID(prefix))
	middlewares = append(middlewares, middleware.Authentication)
	middlewares = append(middlewares, middleware.ContextLogger(cfg.LoggerFactory))
	if cfg.NewrelicApp != nil {
		middlewares = append(middlewares, middleware.Newrelic(cfg.NewrelicApp))
	}
	middlewares = append(middlewares, middleware.Recoverer(cfg.PanicHandler))
	middlewares = append(middlewares, middleware.Prometheus())
	mux := chi.NewMux()
	mux.Use(middlewares...)

	if !cfg.AccessLog.Disable {
		mux.Use(middleware.AccessLog(middleware.AccessLogOptions{}))
	}

	return mux, nil
}
