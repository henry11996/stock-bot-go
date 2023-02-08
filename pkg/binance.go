package pkg

import (
	"os"

	"github.com/adshao/go-binance/v2"
)

func InitBinance() *binance.Client {
	var (
		apiKey    = os.Getenv("BINANCE_API_KEY")
		secretKey = os.Getenv("BINANCE_SECRET_KEY")
	)
	return binance.NewClient(apiKey, secretKey)
}
