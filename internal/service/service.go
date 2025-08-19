package service

import (
	"context"
	"fmt"
	"wb-level0/internal/models"
)

func (s *WBLevel0Service) CreateOrder(ctx context.Context, data []byte) (string, error) {
	order, err := unmarshalToModel[models.Order](data)

	if err != nil {
		fmt.Printf("[service] ошибка парсинга модели %v\n", err)
		return "", err
	}

	id, err := s.repository.CreateOrder(ctx, order)
	if err != nil {
		fmt.Printf("[service] заказ под номером %s не создан %v\n", order.OrderUID, handlePgError(err))
		return "", err
	}

	if err := s.repository.CreateItems(ctx, id, order.Items); err != nil {
		return "", err
	}
	if err := s.repository.CreatePayment(ctx, id, order.Payment); err != nil {
		return "", err
	}
	if err := s.repository.CreateDelivery(ctx, id, order.Delivery); err != nil {
		return "", err
	}

	fmt.Printf("[service] заказ под номером %s создан\n", id)
	return id, nil
}

func (s *WBLevel0Service) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		fmt.Printf("[service] ошибка получения заказа %s %v\n", id, handlePgError(err))
		return models.Order{}, err
	}

	delivery, err := s.repository.GetDeliveryByID(ctx, id)
	if err != nil {
		fmt.Printf("[service] ошибка получения доставки %s %v\n", id, handlePgError(err))
		return models.Order{}, err
	}

	payment, err := s.repository.GetPaymentByID(ctx, id)
	if err != nil {
		fmt.Printf("[service] ошибка получения оплаты %s %v\n", id, handlePgError(err))
		return models.Order{}, err
	}

	items, err := s.repository.GetItemsByID(ctx, id)
	if err != nil {
		fmt.Printf("[service] ошибка получения товаров %s %v\n", id, handlePgError(err))
		return models.Order{}, err
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	fmt.Printf("[service] заказ под номером %s получен\n", id)
	return order, nil
}
