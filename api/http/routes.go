package http

import (
	"iam-service/api/http/handlers"
	"net/http"
)

func NewRouter(authHandler *handlers.AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	return mux
}
