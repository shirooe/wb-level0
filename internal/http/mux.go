package http

import "github.com/gorilla/mux"

// создание роутера
func NewServerMux() *mux.Router {
	return mux.NewRouter()
}
