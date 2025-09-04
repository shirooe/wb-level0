package service

import (
	"context"
	"time"
	"wb-level0/internal/models"

	"go.uber.org/zap"
)

// создание заказа
func (s *WBLevel0Service) CreateOrder(ctx context.Context, data []byte) (string, error) {
	// unmarshal к структуре заказа
	order, err := unmarshalToModel[models.Order](data)
	if err != nil {
		s.log.Info("[service] ошибка парсинга модели", zap.Error(err))
		return "", err
	}

	// валидация структуры
	if err := validate.Struct(order); err != nil {
		s.log.Info("[service] ошибка валидации заказа", zap.Error(err))
		return "", err
	}

	// order_uid
	var id string
	// создание контекста с таймером
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	// запуск транзакции
	err = s.manager.WithTransaction(ctxWithTimeout, func(ctx context.Context) error {
		// ошибка в транзакции
		var errorTx error

		// создание заказа
		id, errorTx = s.repository.CreateOrder(ctxWithTimeout, order)
		if errorTx != nil {
			return errorTx
		}

		// создание доставки
		errorTx = s.repository.CreateDelivery(ctxWithTimeout, id, order.Delivery)
		if errorTx != nil {
			return errorTx
		}

		// создание оплаты
		errorTx = s.repository.CreatePayment(ctxWithTimeout, id, order.Payment)
		if errorTx != nil {
			return errorTx

		}

		// создание предметов
		errorTx = s.repository.CreateItems(ctxWithTimeout, id, order.Items)
		if errorTx != nil {
			return errorTx
		}

		return nil
	})

	if err != nil {
		s.log.Info("[service] ошибка создания заказа", zap.Error(handlePgErrors(err)))
		return "", err
	}

	s.log.Info("[service] заказ создан", zap.String("order_uid", id))
	// сохранение заказа в кэш
	s.cache.Set(id, order)
	s.log.Info("[service] кэш заказов обновлен")
	// возврат order_uid
	return id, nil
}

// поиск заказа по order_uid
func (s *WBLevel0Service) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	// получение данных из кэша
	order, found := s.cache.Get(id)
	if found {
		s.log.Info("[service] заказ получен из кэша", zap.String("order_uid", id))
		return order, nil
	}

	// получение данных из репозитории
	order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		s.log.Info("[service] заказ не получен из базы данных", zap.String("order_uid", id), zap.Error(handlePgErrors(err)))
		return models.Order{}, err
	}

	s.log.Info("[service] заказ получен из базы данных", zap.String("order_uid", id))
	// сохранение данных в кэш
	s.cache.Set(id, order)
	s.log.Info("[service] кэш заказов обновлен")
	// возвращение заказа
	return order, nil
}

// восстановление кэша
func (s *WBLevel0Service) RestoreOrders(ctx context.Context) error {
	// получение всех заказов из репозитории
	orders, err := s.repository.GetAllOrders(ctx)
	if err != nil {
		s.log.Info("[service] ошибка восстановления заказов из базы данных", zap.Error(handlePgErrors(err)))
		return err
	}

	// сохранение заказов в кэш
	s.cache.Restore(orders)
	s.log.Info("[service] кэш заказов восстановлен")
	return nil
}
