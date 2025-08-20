package cache

import (
	"wb-level0/internal/models"
)

func (c *Cache) Set(key string, value models.Order) {
	c.m.Store(key, value)
}

func (c *Cache) Get(key string) (models.Order, bool) {
	value, exist := c.m.Load(key)
	if !exist {
		return models.Order{}, false
	}

	return value.(models.Order), exist
}

func (c *Cache) Restore(orders []models.Order) {
	for _, order := range orders {
		c.m.Store(order.OrderUID, order)
	}
}
