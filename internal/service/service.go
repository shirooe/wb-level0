package service

import (
	"context"
	"wb-level0/internal/models"

	"go.uber.org/zap"
)

func (s *WBLevel0Service) CreateOrder(ctx context.Context, data []byte) (string, error) {
	order, err := unmarshalToModel[models.Order](data)
	if err != nil {
		s.log.Info("[service] ошибка парсинга модели", zap.Error(err))
		return "", err
	}

	id, err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		s.log.Info("[service] заказ не создан", zap.String("order_uid", order.OrderUID), zap.Error(handlePgErrors(err)))
		return "", err
	}

	s.log.Info("[service] заказ создан", zap.String("order_uid", order.OrderUID))
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
