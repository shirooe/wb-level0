package cache

import (
	"sync"

	"go.uber.org/fx"
)

type Cache struct {
	m *sync.Map
}

func NewCache() *Cache {
	return &Cache{
		m: &sync.Map{},
	}
}

func Module() fx.Option {
	return fx.Module("cache",
		fx.Provide(NewCache),
	)
}
