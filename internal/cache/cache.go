package cache

import (
	"wb-level0/internal/models"
)

// сохранение значения
func (c *Cache) Set(key string, value models.Order) {
	c.m.Store(key, value)
}

// получение значения по ключу и маппинг в структуру
func (c *Cache) Get(key string) (models.Order, bool) {
	// загрузка значения из sync.Map
	value, exist := c.m.Load(key)
	if !exist {
		// в случае отсутствия отдаем структуру models.Order
		return models.Order{}, false
	}

	// type assertion значения к структуре models.Order
	return value.(models.Order), exist
}

// сохранение всех поступающих order
// восстановление кэша
func (c *Cache) Restore(orders []models.Order) {
	for _, order := range orders {
		c.m.Store(order.OrderUID, order)
	}
}
