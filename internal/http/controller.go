package http

import (
	"net/http"
	"wb-level0/internal/service"

	"github.com/gorilla/mux"
)

type Controller struct {
	service *service.WBLevel0Service
}

func NewController(service *service.WBLevel0Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/order", c.GetAll).Methods("GET")
	mux.HandleFunc("/order/{id}", c.GetSingle).Methods("GET")
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All OK"))
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetSingle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Single OK"))
	w.WriteHeader(http.StatusOK)
}
