package handler

import "net/http"

type UserHandler interface {
	Create(http.ResponseWriter, *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, userHandler UserHandler) {
	mux.HandleFunc("POST /user", userHandler.Create)
}
