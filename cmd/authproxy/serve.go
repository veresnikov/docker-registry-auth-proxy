package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"

	"github.com/veresnikov/docker-registry-auth-proxy/pkg/application/auth"
	applogger "github.com/veresnikov/docker-registry-auth-proxy/pkg/application/logger"
	infraserver "github.com/veresnikov/docker-registry-auth-proxy/pkg/infrastructure/server"
	"github.com/veresnikov/docker-registry-auth-proxy/pkg/infrastructure/storage"
)

var errServerIsStopped = errors.New("server is stopped")

func serve(cnf *Config, logger applogger.Logger) *cli.Command {
	return &cli.Command{
		Name: "serve",
		Action: func(c *cli.Context) error {
			return serveImpl(c.Context, logger, cnf)
		},
	}
}

func serveImpl(ctx context.Context, logger applogger.Logger, cnf *Config) error {
	hasher := auth.NewPasswordHasher()
	userStorage, err := storage.NewUserStorage(cnf.AccessConfigPath, hasher)
	if err != nil {
		return err
	}
	registryProxy := infraserver.NewRegistryProxy(
		hasher,
		auth.NewService(userStorage),
	)
	err = registryProxy.SetRegistryAddress(cnf.RegistryAddress)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.Use(LoggingMiddleware(logger))
	router.HandleFunc("/v1/_ping", registryProxy.UnsecureServeHTTP)
	router.HandleFunc("/v1/search", registryProxy.UnsecureServeHTTP)
	router.PathPrefix("/").HandlerFunc(registryProxy.SecureServeHTTP)

	server := &http.Server{
		Addr:         cnf.ServeHTTPAddr,
		WriteTimeout: cnf.ServeWriteTimout,
		ReadTimeout:  cnf.ServeReadTimout,
		Handler:      router,
	}

	go func() {
		<-ctx.Done()
		graceCtx, cancel := context.WithTimeout(context.Background(), cnf.ServeGracePeriod)
		defer cancel()
		err = server.Shutdown(graceCtx)
		if err != nil {
			logger.Error(err, "failed server shutdown")
		}
	}()
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return errServerIsStopped
	}
	return err
}

func LoggingMiddleware(logger applogger.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &responseWriterWithStatusCode{
				ResponseWriter: w,
			}
			start := time.Now()
			h.ServeHTTP(wrappedWriter, r)
			duration := time.Since(start)
			logger.WithFields(applogger.Fields{
				"method":   r.Method,
				"url":      r.URL.String(),
				"code":     wrappedWriter.statusCode,
				"duration": duration.String(),
			}).Info("call finished")
		})
	}
}

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (writer *responseWriterWithStatusCode) WriteHeader(statusCode int) {
	writer.ResponseWriter.WriteHeader(statusCode)
	writer.statusCode = statusCode
}
