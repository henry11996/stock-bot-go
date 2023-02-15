package binance

import (
	"os"
	"sync"

	"github.com/adshao/go-binance/v2"
)

var (
	once sync.Once
	c    *binance.Client
)

func GetInstance() *binance.Client {
	once.Do(func() {
		var (
			apiKey    = os.Getenv("BINANCE_API_KEY")
			secretKey = os.Getenv("BINANCE_SECRET_KEY")
		)
		c = binance.NewClient(apiKey, secretKey)
	})
	return c
}
