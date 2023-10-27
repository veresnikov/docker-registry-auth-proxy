package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/veresnikov/docker-registry-auth-proxy/pkg/application/auth"
)

type RegistryProxy interface {
	SetRegistryAddress(address string) error
	SecureServeHTTP(w http.ResponseWriter, r *http.Request)
	UnsecureServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewRegistryProxy(
	passwordHasher auth.PasswordHasher,
	authProvider auth.Service,
) RegistryProxy {
	return &registryProxy{
		passwordHasher: passwordHasher,
		authProvider:   authProvider,
	}
}

type registryProxy struct {
	passwordHasher auth.PasswordHasher
	authProvider   auth.Service
	reverseProxy   *httputil.ReverseProxy
}

func (handler *registryProxy) SetRegistryAddress(address string) error {
	registryURL, err := url.Parse(address)
	if err != nil {
		return err
	}
	handler.reverseProxy = httputil.NewSingleHostReverseProxy(registryURL)
	return nil
}

func (handler *registryProxy) SecureServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler.authMiddleware(r)
	if err != nil {
		handler.errorsHandler(w, err)
		return
	}
	handler.reverseProxy.ServeHTTP(w, r)
}

func (handler *registryProxy) UnsecureServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.reverseProxy.ServeHTTP(w, r)
}

func (handler *registryProxy) authMiddleware(r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if ok {
		passwordHash, err := handler.passwordHasher.Hash(password)
		if err != nil {
			return err
		}
		err = handler.authProvider.Authorize(r.Context(), username, passwordHash)
		if err != nil {
			return err
		}
		return nil
	}
	return auth.ErrUnauthorized
}

func (handler *registryProxy) errorsHandler(w http.ResponseWriter, err error) {
	switch err {
	case nil:
	case auth.ErrUnauthorized:
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
