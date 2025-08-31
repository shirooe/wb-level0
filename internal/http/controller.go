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

	fs := http.FileServer(http.Dir("static"))
	mux.PathPrefix("/").Handler(http.StripPrefix("/", fs))
}

func (c *Controller) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")

	order, err := c.service.GetOrderByID(ctx, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}

	json.NewEncoder(w).Encode(order)
}
