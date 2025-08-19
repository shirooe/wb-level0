package http

import (
	"encoding/json"
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
	mux.HandleFunc("/order/{id}", c.GetOrderByID).Methods("GET")
}
func (c *Controller) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	order, err := c.service.GetOrderByID(ctx, id)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(order)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
