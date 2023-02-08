package pkg

import "github.com/joho/godotenv"

func InitEnv() error {
	return godotenv.Load()
}
