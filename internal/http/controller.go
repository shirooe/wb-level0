package http

import (
	"encoding/json"
	"io"
	"net/http"
	"wb-level0/internal/kafka/producer"
	"wb-level0/internal/service"

	"github.com/gorilla/mux"
)

type Controller struct {
	service *service.WBLevel0Service

	writer             *producer.Producer
	hasTestMessageSent bool
}

func NewController(service *service.WBLevel0Service, writer *producer.Producer) *Controller {
	return &Controller{
		service:            service,
		writer:             writer,
		hasTestMessageSent: false,
	}
}

func (c *Controller) RegisterRoutes(mux *mux.Router) {
	// получение заказа по id
	mux.HandleFunc("/order/{id}", c.GetOrderByID).Methods("GET")

	// для отправки сообщении в кафку
	mux.HandleFunc("/order", c.CreateOrder).Methods("POST")

	// получение статических файлов
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

func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: повторный запуск приложения
	if c.hasTestMessageSent {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Тестовое сообщение уже было отправлено",
			},
		)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}
	defer r.Body.Close()

	if err := c.writer.WriteTestMessage(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}
	c.hasTestMessageSent = true
	w.WriteHeader(http.StatusOK)
}
