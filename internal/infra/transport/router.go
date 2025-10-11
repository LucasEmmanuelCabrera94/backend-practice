package transport

import (
	"backend-practice/internal/infra/transport/handler"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	return mux
}