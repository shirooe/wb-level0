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
	// проверка на отправленность тестовой модели продьюсером
	hasTestMessageSent bool
}

// создание контроллера
func NewController(service *service.WBLevel0Service, writer *producer.Producer) *Controller {
	return &Controller{
		service:            service,
		writer:             writer,
		hasTestMessageSent: false,
	}
}

// регистрация роутов
func (c *Controller) RegisterRoutes(mux *mux.Router) {
	// получение заказа по id
	mux.HandleFunc("/order/{id}", c.GetOrderByID).Methods("GET")

	// для отправки сообщении в кафку
	mux.HandleFunc("/order", c.CreateOrder).Methods("POST")

	// получение статических файлов
	fs := http.FileServer(http.Dir("static"))
	mux.PathPrefix("/").Handler(http.StripPrefix("/", fs))
}

// получения заказа по ID
func (c *Controller) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	// получение контекста
	ctx := r.Context()
	// получение параметра из строки
	vars := mux.Vars(r)
	id := vars["id"]
	// сохранение response content-type
	w.Header().Set("Content-Type", "application/json")

	// получения заказа по ID
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

	// возвращение клиенту модели Order
	json.NewEncoder(w).Encode(order)
}

// создание тестового заказа (продьюсер отправляет сообщение в топик)
func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// настройка content-type
	w.Header().Set("Content-Type", "application/json")

	// TODO: повторный запуск приложения
	// проверка на уже отправленное сообщение
	if c.hasTestMessageSent {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": "Тестовое сообщение уже было отправлено",
			},
		)
		return
	}

	// считывание данных из body (приходит тестовая модель заказа)
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
	// закрытие body
	defer r.Body.Close()

	// продьюсер отправляет тестовую модель в топик
	if err := c.writer.WriteTestMessage(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)
		return
	}
	// сообщение уже ранее отправлено
	c.hasTestMessageSent = true
	w.WriteHeader(http.StatusOK)
}
