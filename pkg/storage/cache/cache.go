package cache

import (
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"
)

var (
	once sync.Once
	c    *goCache.Cache
)

func GetInstance() *goCache.Cache {
	once.Do(func() {
		c = goCache.New(0*time.Minute, 0*time.Minute)
	})
	return c
}
