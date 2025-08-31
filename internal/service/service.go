package service

import (
	"context"
	"time"
	"wb-level0/internal/models"

	"go.uber.org/zap"
)

func (s *WBLevel0Service) CreateOrder(ctx context.Context, data []byte) (string, error) {
	order, err := unmarshalToModel[models.Order](data)
	if err != nil {
		s.log.Info("[service] ошибка парсинга модели", zap.Error(err))
		return "", err
	}

	if err := validate.Struct(order); err != nil {
		s.log.Info("[service] ошибка валидации заказа", zap.Error(err))
		return "", err
	}

	var id string
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	err = s.manager.WithTransaction(ctxWithTimeout, func(ctx context.Context) error {
		var errorTx error

		id, errorTx = s.repository.CreateOrder(ctxWithTimeout, order)
		if errorTx != nil {
			return errorTx
		}

		errorTx = s.repository.CreateDelivery(ctxWithTimeout, id, order.Delivery)
		if errorTx != nil {
			return errorTx
		}

		errorTx = s.repository.CreatePayment(ctxWithTimeout, id, order.Payment)
		if errorTx != nil {
			return errorTx

		}

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
	s.cache.Set(id, order)
	s.log.Info("[service] кэш заказов обновлен")
	return id, nil
}

func (s *WBLevel0Service) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	order, found := s.cache.Get(id)
	if found {
		s.log.Info("[service] заказ получен из кэша", zap.String("order_uid", id))
		return order, nil
	}

	order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		s.log.Info("[service] заказ не получен из базы данных", zap.String("order_uid", id), zap.Error(handlePgErrors(err)))
		return models.Order{}, err
	}

	s.log.Info("[service] заказ получен из базы данных", zap.String("order_uid", id))
	s.cache.Set(id, order)
	s.log.Info("[service] кэш заказов обновлен")
	return order, nil
}

func (s *WBLevel0Service) RestoreOrders(ctx context.Context) error {
	orders, err := s.repository.GetAllOrders(ctx)
	if err != nil {
		s.log.Info("[service] ошибка восстановления заказов из базы данных", zap.Error(handlePgErrors(err)))
		return err
	}

	s.cache.Restore(orders)
	s.log.Info("[service] кэш заказов восстановлен")
	return nil
}
