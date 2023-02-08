package pkg

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func InitCache() *cache.Cache {
	return cache.New(0*time.Minute, 0*time.Minute)
}
