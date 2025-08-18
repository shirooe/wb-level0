package http

import "github.com/gorilla/mux"

func NewServerMux() *mux.Router {
	return mux.NewRouter()
}
