package controllers

import (
	"net/http"
	"sync/atomic"
)

type ServiceController struct {
}

func NewServiceController() *ServiceController {
	return &ServiceController{}
}

func (c *ServiceController) Liveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *ServiceController) Readiness(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
